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
	"github.com/jackc/pgx/v5/pgtype"
)

// TagService handles all business logic for tags
type TagService struct {
	db     *database.DB
	logger *logging.Logger
}

// NewTagService creates a new TagService instance
func NewTagService(db *database.DB, logger *logging.Logger) *TagService {
	return &TagService{
		db:     db,
		logger: logger.WithComponent("tag_service"),
	}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(ctx context.Context, req *models.TagRequest) (*models.Tag, error) {
	// Validate request
	if err := s.validateTagRequest(req); err != nil {
		s.logger.LogError(ctx, err, "Tag validation failed", "tag_name", req.Name)
		return nil, err
	}

	s.logger.Info("Creating new tag", "tag_name", req.Name, "color", req.Color)

	var tag *models.Tag

	// Start write transaction to create tag
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		sqlcTag, err := qtx.CreateTag(ctx, store.CreateTagParams{
			Name:        req.Name,
			Color:       stringToPgText(&req.Color),
			Description: stringToPgText(req.Description),
		})
		if err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}

		tag = s.sqlcToModel(sqlcTag)
		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to create tag", "tag_name", req.Name, "color", req.Color)
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	s.logger.Info("Tag created successfully", "tag_id", tag.ID, "tag_name", tag.Name, "color", tag.Color)
	return tag, nil
}

// GetTag retrieves a single tag by ID
func (s *TagService) GetTag(ctx context.Context, tagID string) (*models.Tag, error) {
	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Getting tag", "tag_id", tagID)

	var tag *models.Tag

	// Read operation to get tag
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTag, err := qtx.GetTagByID(ctx, tagUUID)
		if err != nil {
			return fmt.Errorf("failed to get tag: %w", err)
		}

		tag = s.sqlcToModel(sqlcTag)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Info("Tag not found", "tag_id", tagID)
		} else {
			s.logger.LogError(ctx, err, "Failed to get tag", "tag_id", tagID)
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	s.logger.Info("Tag retrieved successfully", "tag_id", tagID, "tag_name", tag.Name)
	return tag, nil
}

// GetTagByName retrieves a single tag by name
func (s *TagService) GetTagByName(ctx context.Context, name string) (*models.Tag, error) {
	s.logger.Info("Getting tag by name", "tag_name", name)

	var tag *models.Tag

	// Read operation to get tag by name
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTag, err := qtx.GetTagByName(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to get tag by name: %w", err)
		}

		tag = s.sqlcToModel(sqlcTag)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Info("Tag not found by name", "tag_name", name)
		} else {
			s.logger.LogError(ctx, err, "Failed to get tag by name", "tag_name", name)
		}
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}

	s.logger.Info("Tag retrieved successfully by name", "tag_name", name, "tag_id", tag.ID)
	return tag, nil
}

// GetAllTags retrieves all tags
func (s *TagService) GetAllTags(ctx context.Context) ([]*models.Tag, error) {
	s.logger.Info("Getting all tags")

	var tags []*models.Tag

	// Read operation to get all tags
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTags, err := qtx.GetAllTags(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all tags: %w", err)
		}

		tags = make([]*models.Tag, len(sqlcTags))
		for i, sqlcTag := range sqlcTags {
			tags[i] = s.sqlcToModel(sqlcTag)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get all tags")
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}

	s.logger.Info("All tags retrieved successfully", "tags_count", len(tags))
	return tags, nil
}

// GetPopularTags retrieves the most popular tags by usage count
func (s *TagService) GetPopularTags(ctx context.Context, limit int32) ([]*models.Tag, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	s.logger.Info("Getting popular tags", "limit", limit)

	var tags []*models.Tag

	// Read operation to get popular tags
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTags, err := qtx.GetPopularTags(ctx, limit)
		if err != nil {
			return fmt.Errorf("failed to get popular tags: %w", err)
		}

		tags = make([]*models.Tag, len(sqlcTags))
		for i, sqlcTag := range sqlcTags {
			tags[i] = s.sqlcToModel(sqlcTag)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get popular tags", "limit", limit)
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}

	s.logger.Info("Popular tags retrieved successfully", "limit", limit, "returned_count", len(tags))
	return tags, nil
}

// GetRecentlyUsedTags retrieves recently used tags for a user
func (s *TagService) GetRecentlyUsedTags(ctx context.Context, userID string, limit int32) ([]*models.Tag, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	s.logger.Info("Getting recently used tags", "user_id", userID, "limit", limit)

	var tags []*models.Tag

	// Read operation to get recently used tags
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		// Use current time as the upper bound for recently used tags
		sqlcTags, err := qtx.GetRecentlyUsedTags(ctx, store.GetRecentlyUsedTagsParams{
			UserID:    userUUID,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to get recently used tags: %w", err)
		}

		// Limit the results manually since the query doesn't support limit parameter
		actualLimit := len(sqlcTags)
		if int32(actualLimit) > limit {
			actualLimit = int(limit)
		}

		tags = make([]*models.Tag, actualLimit)
		for i := 0; i < actualLimit; i++ {
			tags[i] = s.sqlcRecentlyUsedToModel(sqlcTags[i])
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get recently used tags", "user_id", userID, "limit", limit)
		return nil, fmt.Errorf("failed to get recently used tags: %w", err)
	}

	s.logger.Info("Recently used tags retrieved successfully", "user_id", userID, "limit", limit, "returned_count", len(tags))
	return tags, nil
}

// SearchTags searches for tags by name
func (s *TagService) SearchTags(ctx context.Context, query string, limit int32) ([]*models.Tag, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}

	s.logger.Info("Searching tags", "query", query, "limit", limit)

	var tags []*models.Tag

	// Read operation to search tags
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTags, err := qtx.SearchTags(ctx, store.SearchTagsParams{
			Column1: pgtype.Text{String: query, Valid: true},
			Limit:   limit,
		})
		if err != nil {
			return fmt.Errorf("failed to search tags: %w", err)
		}

		tags = make([]*models.Tag, len(sqlcTags))
		for i, sqlcTag := range sqlcTags {
			tags[i] = s.sqlcToModel(sqlcTag)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to search tags", "query", query, "limit", limit)
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}

	s.logger.Info("Tag search completed successfully", "query", query, "limit", limit, "results_count", len(tags))
	return tags, nil
}

// UpdateTag updates an existing tag
func (s *TagService) UpdateTag(ctx context.Context, tagID string, req *models.TagRequest) (*models.Tag, error) {
	// Validate request
	if err := s.validateTagRequest(req); err != nil {
		s.logger.Warn("Invalid tag update request", "tag_id", tagID, "tag_name", req.Name, "error", err.Error())
		return nil, err
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Updating tag", "tag_id", tagID, "tag_name", req.Name, "color", req.Color)

	var tag *models.Tag

	// Start write transaction to update tag
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		sqlcTag, err := qtx.UpdateTag(ctx, store.UpdateTagParams{
			ID:          tagUUID,
			Name:        req.Name,
			Color:       stringToPgText(&req.Color),
			Description: stringToPgText(req.Description),
		})
		if err != nil {
			return fmt.Errorf("failed to update tag: %w", err)
		}

		tag = s.sqlcToModel(sqlcTag)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Tag not found for update", "tag_id", tagID)
		} else {
			s.logger.LogError(ctx, err, "Failed to update tag", "tag_id", tagID, "tag_name", req.Name)
		}
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	s.logger.Info("Tag updated successfully", "tag_id", tagID, "tag_name", tag.Name, "color", tag.Color)
	return tag, nil
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(ctx context.Context, tagID string) error {
	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Deleting tag", "tag_id", tagID)

	// First verify the tag exists
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		_, err := qtx.GetTagByID(ctx, tagUUID)
		if err != nil {
			return fmt.Errorf("tag not found: %w", err)
		}
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Tag not found for deletion", "tag_id", tagID)
		} else {
			s.logger.LogError(ctx, err, "Failed to verify tag for deletion", "tag_id", tagID)
		}
		return fmt.Errorf("failed to verify tag: %w", err)
	}

	// Now delete the tag
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.DeleteTag(ctx, tagUUID)
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to delete tag", "tag_id", tagID)
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	s.logger.Info("Tag deleted successfully", "tag_id", tagID)
	return nil
}

// GetUserTagUsage retrieves tag usage statistics for a specific user
func (s *TagService) GetUserTagUsage(ctx context.Context, userID string) ([]*models.TagUsage, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting user tag usage", "user_id", userID)

	var tagUsages []*models.TagUsage

	// Read operation to get user tag usage
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		usage, err := qtx.GetUserTagUsage(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get user tag usage: %w", err)
		}

		tagUsages = make([]*models.TagUsage, len(usage))
		for i, u := range usage {
			tagUsages[i] = &models.TagUsage{
				Tag: &models.Tag{
					ID:          u.ID,
					Name:        u.Name,
					Color:       pgTextToStringRequired(u.Color),
					Description: pgTextToString(u.Description),
					UsageCount:  0,           // This is global usage count, not user-specific
					CreatedAt:   time.Time{}, // Not available in this query
				},
				UsageCount: int(u.UserUsageCount),
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get user tag usage", "user_id", userID)
		return nil, fmt.Errorf("failed to get user tag usage: %w", err)
	}

	s.logger.Info("User tag usage retrieved successfully", "user_id", userID, "usage_entries_count", len(tagUsages))
	return tagUsages, nil
}

// GetLogEntriesForTag retrieves all log entries that have a specific tag
func (s *TagService) GetLogEntriesForTag(ctx context.Context, tagID string) ([]*models.LogEntry, error) {
	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return nil, fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Getting log entries for tag", "tag_id", tagID)

	var logEntries []*models.LogEntry

	// Read operation to get log entries for tag
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcEntries, err := qtx.GetLogEntriesForTag(ctx, tagUUID)
		if err != nil {
			return fmt.Errorf("failed to get log entries for tag: %w", err)
		}

		logEntries = make([]*models.LogEntry, len(sqlcEntries))
		for i, sqlcEntry := range sqlcEntries {
			logEntries[i] = s.sqlcLogEntryToModel(sqlcEntry)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get log entries for tag", "tag_id", tagID)
		return nil, fmt.Errorf("failed to get log entries for tag: %w", err)
	}

	s.logger.Info("Log entries for tag retrieved successfully", "tag_id", tagID, "entries_count", len(logEntries))
	return logEntries, nil
}

// AddTagToLogEntry associates a tag with a log entry
func (s *TagService) AddTagToLogEntry(ctx context.Context, logEntryID, tagID string) error {
	logEntryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format", "log_entry_id", logEntryID)
		return fmt.Errorf("invalid log entry ID: %w", err)
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Adding tag to log entry", "log_entry_id", logEntryID, "tag_id", tagID)

	// Start write transaction to add tag to log entry
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.AddTagToLogEntry(ctx, store.AddTagToLogEntryParams{
			LogEntryID: logEntryUUID,
			TagID:      tagUUID,
		})
	}); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			s.logger.Warn("Tag already associated with log entry", "log_entry_id", logEntryID, "tag_id", tagID)
		} else {
			s.logger.LogError(ctx, err, "Failed to add tag to log entry", "log_entry_id", logEntryID, "tag_id", tagID)
		}
		return fmt.Errorf("failed to add tag to log entry: %w", err)
	}

	s.logger.Info("Tag added to log entry successfully", "log_entry_id", logEntryID, "tag_id", tagID)
	return nil
}

// RemoveTagFromLogEntry removes a tag association from a log entry
func (s *TagService) RemoveTagFromLogEntry(ctx context.Context, logEntryID, tagID string) error {
	logEntryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format", "log_entry_id", logEntryID)
		return fmt.Errorf("invalid log entry ID: %w", err)
	}

	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid tag ID format", "tag_id", tagID)
		return fmt.Errorf("invalid tag ID: %w", err)
	}

	s.logger.Info("Removing tag from log entry", "log_entry_id", logEntryID, "tag_id", tagID)

	// Start write transaction to remove tag from log entry
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.RemoveTagFromLogEntry(ctx, store.RemoveTagFromLogEntryParams{
			LogEntryID: logEntryUUID,
			TagID:      tagUUID,
		})
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Tag association not found for removal", "log_entry_id", logEntryID, "tag_id", tagID)
		} else {
			s.logger.LogError(ctx, err, "Failed to remove tag from log entry", "log_entry_id", logEntryID, "tag_id", tagID)
		}
		return fmt.Errorf("failed to remove tag from log entry: %w", err)
	}

	s.logger.Info("Tag removed from log entry successfully", "log_entry_id", logEntryID, "tag_id", tagID)
	return nil
}

// GetTagsForLogEntry retrieves all tags associated with a log entry
func (s *TagService) GetTagsForLogEntry(ctx context.Context, logEntryID string) ([]*models.Tag, error) {
	logEntryUUID, err := uuid.Parse(logEntryID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid log entry ID format", "log_entry_id", logEntryID)
		return nil, fmt.Errorf("invalid log entry ID: %w", err)
	}

	s.logger.Info("Getting tags for log entry", "log_entry_id", logEntryID)

	var tags []*models.Tag

	// Read operation to get tags for log entry
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcTags, err := qtx.GetTagsForLogEntry(ctx, logEntryUUID)
		if err != nil {
			return fmt.Errorf("failed to get tags for log entry: %w", err)
		}

		tags = make([]*models.Tag, len(sqlcTags))
		for i, sqlcTag := range sqlcTags {
			tags[i] = s.sqlcToModel(sqlcTag)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get tags for log entry", "log_entry_id", logEntryID)
		return nil, fmt.Errorf("failed to get tags for log entry: %w", err)
	}

	s.logger.Info("Tags for log entry retrieved successfully", "log_entry_id", logEntryID, "tags_count", len(tags))
	return tags, nil
}

// CleanupUnusedTags removes tags that are not associated with any log entries
func (s *TagService) CleanupUnusedTags(ctx context.Context) error {
	s.logger.Info("Starting cleanup of unused tags")

	// Start write transaction to cleanup unused tags
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.CleanupUnusedTags(ctx)
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to cleanup unused tags")
		return fmt.Errorf("failed to cleanup unused tags: %w", err)
	}

	s.logger.Info("Unused tags cleanup completed successfully")
	return nil
}

// EnsureTagExists creates a tag if it doesn't exist, or returns the existing one
func (s *TagService) EnsureTagExists(ctx context.Context, name, color string) (*models.Tag, error) {
	s.logger.Info("Ensuring tag exists", "tag_name", name, "color", color)

	// First try to get existing tag
	tag, err := s.GetTagByName(ctx, name)
	if err == nil {
		s.logger.Info("Tag already exists", "tag_name", name, "tag_id", tag.ID)
		return tag, nil
	}

	s.logger.Info("Tag not found, creating new tag", "tag_name", name, "color", color)

	// If tag doesn't exist, create it
	req := &models.TagRequest{
		Name:        name,
		Color:       color,
		Description: nil,
	}

	createdTag, err := s.CreateTag(ctx, req)
	if err != nil {
		s.logger.LogError(ctx, err, "Failed to create tag in EnsureTagExists", "tag_name", name, "color", color)
		return nil, err
	}

	s.logger.Info("Tag created in EnsureTagExists", "tag_name", name, "tag_id", createdTag.ID)
	return createdTag, nil
}

// validateTagRequest validates the tag request
func (s *TagService) validateTagRequest(req *models.TagRequest) error {
	if req.Name == "" {
		return fmt.Errorf("tag name is required")
	}

	if len(req.Name) > 100 {
		return fmt.Errorf("tag name must be at most 100 characters")
	}

	if req.Color == "" {
		return fmt.Errorf("tag color is required")
	}

	if err := models.ValidateHexColor(req.Color); err != nil {
		return fmt.Errorf("invalid color format: %w", err)
	}

	if req.Description != nil && len(*req.Description) > 500 {
		return fmt.Errorf("tag description must be at most 500 characters")
	}

	return nil
}

// sqlcToModel converts SQLC Tag to models.Tag
func (s *TagService) sqlcToModel(sqlcTag store.Tag) *models.Tag {
	return &models.Tag{
		ID:          sqlcTag.ID,
		Name:        sqlcTag.Name,
		Color:       pgTextToStringRequired(sqlcTag.Color),
		Description: pgTextToString(sqlcTag.Description),
		UsageCount:  int(pgInt4ToInt32(sqlcTag.UsageCount)),
		CreatedAt:   pgTimestamptzToTime(sqlcTag.CreatedAt),
	}
}

// sqlcRecentlyUsedToModel converts SQLC GetRecentlyUsedTagsRow to models.Tag
func (s *TagService) sqlcRecentlyUsedToModel(sqlcRow store.GetRecentlyUsedTagsRow) *models.Tag {
	return &models.Tag{
		ID:          sqlcRow.ID,
		Name:        sqlcRow.Name,
		Color:       pgTextToStringRequired(sqlcRow.Color),
		Description: pgTextToString(sqlcRow.Description),
		UsageCount:  int(pgInt4ToInt32(sqlcRow.UsageCount)),
		CreatedAt:   pgTimestamptzToTime(sqlcRow.CreatedAt),
	}
}

// sqlcLogEntryToModel converts SQLC LogEntry to models.LogEntry (helper for GetLogEntriesForTag)
func (s *TagService) sqlcLogEntryToModel(sqlcEntry store.LogEntry) *models.LogEntry {
	return &models.LogEntry{
		ID:              sqlcEntry.ID,
		UserID:          sqlcEntry.UserID,
		Title:           sqlcEntry.Title,
		Description:     pgTextToString(sqlcEntry.Description),
		Type:            models.ActivityType(sqlcEntry.Type),
		ProjectID:       pgUUIDToUUIDPtr(sqlcEntry.ProjectID),
		StartTime:       pgTimestamptzToTime(sqlcEntry.StartTime),
		EndTime:         pgTimestamptzToTime(sqlcEntry.EndTime),
		DurationMinutes: int(pgInt4ToInt32(sqlcEntry.DurationMinutes)),
		ValueRating:     models.ValueRating(sqlcEntry.ValueRating),
		ImpactLevel:     models.ImpactLevel(sqlcEntry.ImpactLevel),
		CreatedAt:       pgTimestamptzToTime(sqlcEntry.CreatedAt),
		UpdatedAt:       pgTimestamptzToTime(sqlcEntry.UpdatedAt),
	}
}
