package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/google/uuid"
)

// LogEntryService handles all business logic for activity log entries
type LogEntryService struct {
	db     *database.DB
	logger *logging.Logger
}

// NewLogEntryService creates a new LogEntryService instance
func NewLogEntryService(db *database.DB, logger *logging.Logger) *LogEntryService {
	return &LogEntryService{
		db:     db,
		logger: logger.WithComponent("log_entry_service"),
	}
}

// CreateLogEntry creates a new log entry with associated tags
func (s *LogEntryService) CreateLogEntry(ctx context.Context, userID string, req *models.LogEntryRequest) (*models.LogEntry, error) {
	// Validate request
	if err := s.validateLogEntryRequest(req); err != nil {
		s.logger.LogError(ctx, err, "Log entry validation failed", "user_id", userID, "type", req.Type)
		return nil, err
	}

	// Convert userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in CreateLogEntry", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	duration := int(req.EndTime.Sub(req.StartTime).Minutes())
	s.logger.Info("Creating new log entry",
		"user_id", userID,
		"type", req.Type,
		"duration_minutes", duration,
		"project_id", req.ProjectID,
		"title", req.Title)

	var logEntry *models.LogEntry

	// Start write transaction to add log entry and tags
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		// Create log entry
		sqlcEntry, err := qtx.CreateLogEntry(ctx, store.CreateLogEntryParams{
			UserID:      userUUID,
			ProjectID:   uuidToPgUUID(req.ProjectID),
			Title:       req.Title,
			Description: stringToPgText(req.Description),
			Type:        string(req.Type),
			StartTime:   timeToPgTimestamptz(req.StartTime),
			EndTime:     timeToPgTimestamptz(req.EndTime),
			ValueRating: string(req.ValueRating),
			ImpactLevel: string(req.ImpactLevel),
		})
		if err != nil {
			s.logger.LogError(ctx, err, "Database error creating log entry", "user_id", userID, "title", req.Title)
			return fmt.Errorf("failed to create log entry: %w", err)
		}

		// Handle tags if provided
		if len(req.Tags) > 0 {
			s.logger.Info("Adding tags to log entry", "log_entry_id", sqlcEntry.ID, "tags_count", len(req.Tags), "tags", req.Tags)
			for _, tagName := range req.Tags {
				tagID, err := s.ensureTagExists(ctx, qtx, tagName)
				if err != nil {
					s.logger.LogError(ctx, err, "Failed to handle tag", "tag_name", tagName, "log_entry_id", sqlcEntry.ID)
					return fmt.Errorf("failed to handle tag %s: %w", tagName, err)
				}

				err = qtx.AddTagToLogEntry(ctx, store.AddTagToLogEntryParams{
					LogEntryID: sqlcEntry.ID,
					TagID:      tagID,
				})
				if err != nil {
					s.logger.LogError(ctx, err, "Failed to associate tag", "tag_name", tagName, "log_entry_id", sqlcEntry.ID, "tag_id", tagID)
					return fmt.Errorf("failed to associate tag: %w", err)
				}
			}
		}

		// Convert to model and return
		logEntry = s.sqlcToModel(sqlcEntry)
		logEntry.Tags = req.Tags

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Transaction failed for log entry creation", "user_id", userID, "title", req.Title)
		return nil, fmt.Errorf("failed to create log entry: %w", err)
	}

	s.logger.Info("Log entry created successfully",
		"user_id", userID,
		"log_entry_id", logEntry.ID,
		"title", logEntry.Title,
		"duration_minutes", duration,
		"tags_count", len(req.Tags))

	return logEntry, nil
}

// GetLogEntry retrieves a single log entry by ID
func (s *LogEntryService) GetLogEntry(ctx context.Context, userID, logEntryID string) (*models.LogEntry, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetLogEntry", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	entryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format", "user_id", userID, "log_entry_id", logEntryID)
		return nil, fmt.Errorf("invalid log entry ID: %w", err)
	}

	s.logger.Info("Getting log entry", "user_id", userID, "log_entry_id", logEntryID)

	var logEntry *models.LogEntry

	// Read operation to get log entry
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		// Get log entry
		sqlcEntry, err := qtx.GetLogEntryByID(ctx, entryUUID)
		if err != nil {
			return fmt.Errorf("failed to get log entry: %w", err)
		}

		// Verify ownership
		if sqlcEntry.UserID != userUUID {
			s.logger.Warn("Unauthorized access attempt to log entry", "user_id", userID, "log_entry_id", logEntryID, "owner_id", sqlcEntry.UserID)
			return fmt.Errorf("log entry not found")
		}

		// Get associated tags
		tags, err := qtx.GetTagsForLogEntry(ctx, entryUUID)
		if err != nil {
			return fmt.Errorf("failed to get tags: %w", err)
		}

		logEntry = s.sqlcToModel(sqlcEntry)
		logEntry.Tags = make([]string, len(tags))
		for i, tag := range tags {
			logEntry.Tags[i] = tag.Name
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get log entry", "user_id", userID, "log_entry_id", logEntryID)
		return nil, fmt.Errorf("failed to get log entry: %w", err)
	}

	s.logger.Info("Log entry retrieved successfully",
		"user_id", userID,
		"log_entry_id", logEntryID,
		"title", logEntry.Title,
		"tags_count", len(logEntry.Tags))

	return logEntry, nil
}

// LogEntryFilters defines filtering options for log entries
type LogEntryFilters struct {
	StartDate   time.Time
	EndDate     time.Time
	ProjectID   *string
	Type        *models.ActivityType
	ValueRating *models.ValueRating
	ImpactLevel *models.ImpactLevel
	Tags        []string
	Limit       int32 // for pagination
	Offset      int32 // for pagination
}

// GetLogEntries retrieves log entries with optional filtering
func (s *LogEntryService) GetLogEntries(ctx context.Context, userID string, filters *LogEntryFilters) ([]*models.LogEntry, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetLogEntries", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting log entries", "user_id", userID, "has_filters", filters != nil)

	var (
		sqlcEntries []store.LogEntry
		entries     []*models.LogEntry
	)

	// Start read transaction to get log entries
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		if filters != nil && (!filters.StartDate.IsZero() || !filters.EndDate.IsZero()) {
			// Use date range query
			startDate := filters.StartDate
			endDate := filters.EndDate
			if startDate.IsZero() {
				startDate = time.Now().AddDate(-1, 0, 0) // Default to 1 year ago
			}
			if endDate.IsZero() {
				endDate = time.Now()
			}

			sqlcEntries, err = qtx.GetLogEntriesByUserAndDateRange(ctx, store.GetLogEntriesByUserAndDateRangeParams{
				UserID:    userUUID,
				StartTime: timeToPgTimestamptz(startDate),
				EndTime:   timeToPgTimestamptz(endDate),
			})
		} else {
			// Get all entries for user with pagination support
			limit := int32(100) // default limit
			offset := int32(0)  // default offset

			if filters != nil {
				if filters.Limit > 0 {
					limit = filters.Limit
				}
				if filters.Offset > 0 {
					offset = filters.Offset
				}
			}

			sqlcEntries, err = qtx.GetLogEntriesByUser(ctx, store.GetLogEntriesByUserParams{
				UserID: userUUID,
				Limit:  limit,
				Offset: offset,
			})
		}
		if err != nil {
			return err
		}

		// Convert to models
		entries = make([]*models.LogEntry, len(sqlcEntries))
		for i, sqlcEntry := range sqlcEntries {
			entries[i] = s.sqlcToModel(sqlcEntry)

			// Get tags for each entry
			tags, err := qtx.GetTagsForLogEntry(ctx, sqlcEntry.ID)
			if err == nil {
				entries[i].Tags = make([]string, len(tags))
				for j, tag := range tags {
					entries[i].Tags[j] = tag.Name
				}
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get log entries", "user_id", userID)
		return nil, fmt.Errorf("failed to get log entries: %w", err)
	}

	// Apply additional filters
	if filters != nil {
		entries = s.applyFilters(entries, filters)
	}

	s.logger.Info("Log entries retrieved successfully",
		"user_id", userID,
		"entries_count", len(entries),
		"filters_applied", filters != nil)

	return entries, nil
}

// UpdateLogEntry updates an existing log entry
func (s *LogEntryService) UpdateLogEntry(ctx context.Context, userID, logEntryID string, req *models.LogEntryRequest) (*models.LogEntry, error) {
	// Validate request
	if err := s.validateLogEntryRequest(req); err != nil {
		s.logger.LogError(ctx, err, "Log entry validation failed for update", "user_id", userID, "log_entry_id", logEntryID)
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in UpdateLogEntry", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	entryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format for update", "user_id", userID, "log_entry_id", logEntryID)
		return nil, fmt.Errorf("invalid log entry ID: %w", err)
	}

	duration := int(req.EndTime.Sub(req.StartTime).Minutes())
	s.logger.Info("Updating log entry",
		"user_id", userID,
		"log_entry_id", logEntryID,
		"title", req.Title,
		"duration_minutes", duration)

	var logEntry *models.LogEntry

	// Start write transaction to add log entry and tags
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		// Update log entry
		sqlcEntry, err := qtx.UpdateLogEntry(ctx, store.UpdateLogEntryParams{
			ID:          entryUUID,
			Title:       req.Title,
			Description: stringToPgText(req.Description),
			Type:        string(req.Type),
			ProjectID:   uuidToPgUUID(req.ProjectID),
			StartTime:   timeToPgTimestamptz(req.StartTime),
			EndTime:     timeToPgTimestamptz(req.EndTime),
			ValueRating: string(req.ValueRating),
			ImpactLevel: string(req.ImpactLevel),
			UserID:      userUUID,
		})
		if err != nil {
			return err
		}

		// Update tags - remove old associations and create new ones
		// First, get existing tags to remove them
		existingTags, err := qtx.GetTagsForLogEntry(ctx, entryUUID)
		if err != nil {
			return fmt.Errorf("failed to get existing tags: %w", err)
		}

		// Remove existing tag associations
		for _, tag := range existingTags {
			err = qtx.RemoveTagFromLogEntry(ctx, store.RemoveTagFromLogEntryParams{
				LogEntryID: entryUUID,
				TagID:      tag.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to remove existing tag: %w", err)
			}
		}

		if len(req.Tags) > 0 {
			for _, tagName := range req.Tags {
				tagID, err := s.ensureTagExists(ctx, qtx, tagName)
				if err != nil {
					return fmt.Errorf("failed to handle tag %s: %w", tagName, err)
				}

				err = qtx.AddTagToLogEntry(ctx, store.AddTagToLogEntryParams{
					LogEntryID: entryUUID,
					TagID:      tagID,
				})
				if err != nil {
					return fmt.Errorf("failed to associate tag: %w", err)
				}
			}
		}

		// Convert to model
		logEntry = s.sqlcToModel(sqlcEntry)
		logEntry.Tags = req.Tags

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Transaction failed for log entry update", "user_id", userID, "log_entry_id", logEntryID)
		return nil, fmt.Errorf("failed to update log entry: %w", err)
	}

	s.logger.Info("Log entry updated successfully",
		"user_id", userID,
		"log_entry_id", logEntryID,
		"title", logEntry.Title,
		"tags_count", len(req.Tags))

	return logEntry, nil
}

// DeleteLogEntry deletes a log entry
func (s *LogEntryService) DeleteLogEntry(ctx context.Context, userID, logEntryID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in DeleteLogEntry", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	entryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format for delete", "user_id", userID, "log_entry_id", logEntryID)
		return fmt.Errorf("invalid log entry ID: %w", err)
	}

	s.logger.Info("Deleting log entry", "user_id", userID, "log_entry_id", logEntryID)

	// Start write transaction to delete log entry
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		rowsAffected, err := qtx.DeleteLogEntry(ctx, store.DeleteLogEntryParams{
			ID:     entryUUID,
			UserID: userUUID,
		})
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("log entry not found")
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to delete log entry", "user_id", userID, "log_entry_id", logEntryID)
		return fmt.Errorf("failed to delete log entry: %w", err)
	}

	s.logger.Info("Log entry deleted successfully", "user_id", userID, "log_entry_id", logEntryID)
	return nil
}

// validateLogEntryRequest validates the log entry request
func (s *LogEntryService) validateLogEntryRequest(req *models.LogEntryRequest) error {
	// Title validation
	if req.Title == "" {
		s.logger.Warn("Validation failed: title is required")
		return fmt.Errorf("title is required")
	}

	// Check for whitespace-only title
	if strings.TrimSpace(req.Title) == "" {
		s.logger.Warn("Validation failed: title cannot be whitespace only")
		return fmt.Errorf("title cannot be whitespace only")
	}

	if len(req.Title) > 200 {
		s.logger.Warn("Validation failed: title too long", "title_length", len(req.Title))
		return fmt.Errorf("title must be at most 200 characters")
	}

	// Description validation
	if req.Description != nil && len(*req.Description) > 1000 {
		s.logger.Warn("Validation failed: description too long", "description_length", len(*req.Description))
		return fmt.Errorf("description must be at most 1000 characters")
	}

	// Time validation
	if req.StartTime.Equal(req.EndTime) {
		s.logger.Warn("Validation failed: start time equals end time", "start_time", req.StartTime, "end_time", req.EndTime)
		return fmt.Errorf("end time must be after start time")
	}

	if req.EndTime.Before(req.StartTime) {
		s.logger.Warn("Validation failed: end time before start time", "start_time", req.StartTime, "end_time", req.EndTime)
		return fmt.Errorf("end time must be after start time")
	}

	// Activity type validation
	if !req.Type.IsValid() {
		s.logger.Warn("Validation failed: invalid activity type", "type", req.Type)
		return fmt.Errorf("invalid activity type: %s", req.Type)
	}

	// Value rating validation
	if !req.ValueRating.IsValid() {
		s.logger.Warn("Validation failed: invalid value rating", "value_rating", req.ValueRating)
		return fmt.Errorf("invalid value rating: %s", req.ValueRating)
	}

	// Impact level validation
	if !req.ImpactLevel.IsValid() {
		s.logger.Warn("Validation failed: invalid impact level", "impact_level", req.ImpactLevel)
		return fmt.Errorf("invalid impact level: %s", req.ImpactLevel)
	}

	return nil
}

// sqlcToModel converts SQLC LogEntry to models.LogEntry
func (s *LogEntryService) sqlcToModel(sqlcEntry store.LogEntry) *models.LogEntry {
	return &models.LogEntry{
		ID:              sqlcEntry.ID,
		UserID:          sqlcEntry.UserID,
		Title:           sqlcEntry.Title,
		Description:     pgTextToString(sqlcEntry.Description),
		Type:            models.ActivityType(sqlcEntry.Type),
		ProjectID:       pgUUIDToUUID(sqlcEntry.ProjectID),
		StartTime:       pgTimestamptzToTime(sqlcEntry.StartTime),
		EndTime:         pgTimestamptzToTime(sqlcEntry.EndTime),
		DurationMinutes: pgInt4ToInt(sqlcEntry.DurationMinutes),
		ValueRating:     models.ValueRating(sqlcEntry.ValueRating),
		ImpactLevel:     models.ImpactLevel(sqlcEntry.ImpactLevel),
		CreatedAt:       pgTimestamptzToTime(sqlcEntry.CreatedAt),
		UpdatedAt:       pgTimestamptzToTime(sqlcEntry.UpdatedAt),
	}
}

// ensureTagExists gets or creates a tag by name
func (s *LogEntryService) ensureTagExists(ctx context.Context, qtx *store.Queries, tagName string) (uuid.UUID, error) {
	// Try to get existing tag
	tag, err := qtx.GetTagByName(ctx, tagName)
	if err == nil {
		s.logger.Info("Using existing tag", "tag_name", tagName, "tag_id", tag.ID)
		return tag.ID, nil
	}

	// Create new tag if it doesn't exist
	s.logger.Info("Creating new tag", "tag_name", tagName)
	newTag, err := qtx.CreateTag(ctx, store.CreateTagParams{
		Name:  tagName,
		Color: stringToPgText(&[]string{"#3B82F6"}[0]), // Default blue color
	})
	if err != nil {
		s.logger.LogError(ctx, err, "Failed to create new tag", "tag_name", tagName)
		return uuid.Nil, err
	}

	s.logger.Info("New tag created successfully", "tag_name", tagName, "tag_id", newTag.ID)
	return newTag.ID, nil
}

// applyFilters applies additional filters to log entries
func (s *LogEntryService) applyFilters(entries []*models.LogEntry, filters *LogEntryFilters) []*models.LogEntry {
	if filters == nil {
		return entries
	}

	filtered := make([]*models.LogEntry, 0, len(entries))

	for _, entry := range entries {
		// Apply project filter
		if filters.ProjectID != nil {
			if entry.ProjectID == nil || entry.ProjectID.String() != *filters.ProjectID {
				continue
			}
		}

		// Apply type filter
		if filters.Type != nil && entry.Type != *filters.Type {
			continue
		}

		// Apply value rating filter
		if filters.ValueRating != nil && entry.ValueRating != *filters.ValueRating {
			continue
		}

		// Apply impact level filter
		if filters.ImpactLevel != nil && entry.ImpactLevel != *filters.ImpactLevel {
			continue
		}

		// Apply tags filter
		if len(filters.Tags) > 0 {
			hasRequiredTag := false
			for _, requiredTag := range filters.Tags {
				for _, entryTag := range entry.Tags {
					if entryTag == requiredTag {
						hasRequiredTag = true
						break
					}
				}
				if hasRequiredTag {
					break
				}
			}
			if !hasRequiredTag {
				continue
			}
		}

		filtered = append(filtered, entry)
	}

	return filtered
}
