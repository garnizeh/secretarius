package services

import (
	"context"
	"fmt"
	"time"

	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/google/uuid"
)

// AnalyticsService handles all business logic for analytics and productivity metrics
type AnalyticsService struct {
	db     *database.DB
	logger *logging.Logger
}

// NewAnalyticsService creates a new AnalyticsService instance
func NewAnalyticsService(db *database.DB, logger *logging.Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger.WithComponent("analytics_service"),
	}
}

// GetProductivityMetrics retrieves comprehensive productivity metrics for a user
func (s *AnalyticsService) GetProductivityMetrics(ctx context.Context, userID string, startDate, endDate time.Time) (*models.ProductivityMetrics, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetProductivityMetrics", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting productivity metrics", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var metrics *models.ProductivityMetrics

	// Read operation to get productivity metrics
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		// Get activity type distribution
		typeDistribution, err := qtx.GetActivityTypeDistribution(ctx, store.GetActivityTypeDistributionParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get activity type distribution: %w", err)
		}

		// Get value rating distribution
		valueDistribution, err := qtx.GetValueRatingDistribution(ctx, store.GetValueRatingDistributionParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get value rating distribution: %w", err)
		}

		// Get impact level distribution
		impactDistribution, err := qtx.GetImpactLevelDistribution(ctx, store.GetImpactLevelDistributionParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get impact level distribution: %w", err)
		}

		// Process distributions into model format
		activityBreakdown := make(map[models.ActivityType]int)
		totalActivities := 0
		totalMinutes := 0
		highValueActivities := 0

		for _, dist := range typeDistribution {
			activityType := models.ActivityType(dist.Type)
			activityBreakdown[activityType] = int(dist.TotalMinutes)
			totalActivities += int(dist.EntryCount)
			totalMinutes += int(dist.TotalMinutes)
		}

		valueRatingMap := make(map[models.ValueRating]int)
		for _, dist := range valueDistribution {
			valueRating := models.ValueRating(dist.ValueRating)
			valueRatingMap[valueRating] = int(dist.EntryCount)
			if valueRating == models.ValueHigh || valueRating == models.ValueCritical {
				highValueActivities += int(dist.EntryCount)
			}
		}

		impactLevelMap := make(map[models.ImpactLevel]int)
		for _, dist := range impactDistribution {
			impactLevel := models.ImpactLevel(dist.ImpactLevel)
			impactLevelMap[impactLevel] = int(dist.EntryCount)
		}

		// Calculate daily averages
		daysDiff := int(endDate.Sub(startDate).Hours()/24) + 1
		dailyAverages := map[string]float64{
			"activities_per_day": float64(totalActivities) / float64(daysDiff),
			"minutes_per_day":    float64(totalMinutes) / float64(daysDiff),
			"hours_per_day":      float64(totalMinutes) / float64(daysDiff) / 60,
		}

		// Get project count (simplified calculation)
		projectsWorked := 0
		for _, dist := range typeDistribution {
			if dist.EntryCount > 0 {
				projectsWorked++ // This is a simplified count
			}
		}

		metrics = &models.ProductivityMetrics{
			TotalActivities:     totalActivities,
			TotalMinutes:        totalMinutes,
			ProjectsWorked:      projectsWorked,
			HighValueActivities: highValueActivities,
			ActivityBreakdown:   activityBreakdown,
			ValueDistribution:   valueRatingMap,
			ImpactDistribution:  impactLevelMap,
			DailyAverages:       dailyAverages,
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get productivity metrics", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get productivity metrics: %w", err)
	}

	s.logger.Info("Successfully retrieved productivity metrics",
		"user_id", userID,
		"total_activities", metrics.TotalActivities,
		"total_minutes", metrics.TotalMinutes,
		"projects_worked", metrics.ProjectsWorked)

	return metrics, nil
}

// GetActivitySummary retrieves activity summary for a user within a date range
func (s *AnalyticsService) GetActivitySummary(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.ActivitySummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetActivitySummary", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var summaries []*models.ActivitySummary

	// Read operation to get activity summary
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		typeDistribution, err := qtx.GetActivityTypeDistribution(ctx, store.GetActivityTypeDistributionParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get activity type distribution: %w", err)
		}

		summaries = make([]*models.ActivitySummary, len(typeDistribution))
		for i, dist := range typeDistribution {
			summaries[i] = &models.ActivitySummary{
				Type:         models.ActivityType(dist.Type),
				Count:        int(dist.EntryCount),
				TotalMinutes: int(dist.TotalMinutes),
				AvgMinutes:   dist.AvgDuration,
				Week:         startDate, // Using start date as reference
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get activity summary: %w", err)
	}

	s.logger.Info("Successfully retrieved activity summary", "user_id", userID, "summary_count", len(summaries))
	return summaries, nil
}

// GetWeeklyActivitySummary retrieves weekly activity summary for a user
func (s *AnalyticsService) GetWeeklyActivitySummary(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.ActivitySummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetWeeklyActivitySummary", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting weekly activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var summaries []*models.ActivitySummary

	// Read operation to get weekly activity summary
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		weeklyData, err := qtx.GetWeeklyActivitySummary(ctx, store.GetWeeklyActivitySummaryParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get weekly activity summary: %w", err)
		}

		summaries = make([]*models.ActivitySummary, len(weeklyData))
		for i, week := range weeklyData {
			summaries[i] = &models.ActivitySummary{
				Type:         models.ActivityDevelopment, // Default type for weekly summary
				Count:        int(week.EntryCount),
				TotalMinutes: int(week.TotalMinutes),
				AvgMinutes:   week.AvgDuration,
				Week:         startDate, // Using start date as reference since WeekStart is interval
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get weekly activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get weekly activity summary: %w", err)
	}

	s.logger.Info("Successfully retrieved weekly activity summary", "user_id", userID, "weeks_count", len(summaries))
	return summaries, nil
}

// GetMonthlyActivitySummary retrieves monthly activity summary for a user
func (s *AnalyticsService) GetMonthlyActivitySummary(ctx context.Context, userID string, startDate, endDate time.Time) ([]*models.ActivitySummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetMonthlyActivitySummary", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting monthly activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var summaries []*models.ActivitySummary

	// Read operation to get monthly activity summary
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		monthlyData, err := qtx.GetMonthlyActivitySummary(ctx, store.GetMonthlyActivitySummaryParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get monthly activity summary: %w", err)
		}

		summaries = make([]*models.ActivitySummary, len(monthlyData))
		for i, month := range monthlyData {
			summaries[i] = &models.ActivitySummary{
				Type:         models.ActivityDevelopment, // Default type for monthly summary
				Count:        int(month.EntryCount),
				TotalMinutes: int(month.TotalMinutes),
				AvgMinutes:   month.AvgDuration,
				Week:         startDate, // Using start date as reference since MonthStart is interval
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get monthly activity summary", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get monthly activity summary: %w", err)
	}

	s.logger.Info("Successfully retrieved monthly activity summary", "user_id", userID, "months_count", len(summaries))
	return summaries, nil
}

// GetProductivityByDayOfWeek retrieves productivity metrics by day of week
func (s *AnalyticsService) GetProductivityByDayOfWeek(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]float64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetProductivityByDayOfWeek", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting productivity by day of week", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var productivity map[string]float64

	// Read operation to get productivity by day of week
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		weekdayData, err := qtx.GetProductivityByDayOfWeek(ctx, store.GetProductivityByDayOfWeekParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get productivity by day of week: %w", err)
		}

		productivity = make(map[string]float64)
		for _, day := range weekdayData {
			productivity[day.DayName] = float64(day.TotalMinutes)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get productivity by day of week", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get productivity by day of week: %w", err)
	}

	s.logger.Info("Successfully retrieved productivity by day of week", "user_id", userID, "days_count", len(productivity))
	return productivity, nil
}

// GetProductivityByHour retrieves productivity metrics by hour of day
func (s *AnalyticsService) GetProductivityByHour(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]float64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetProductivityByHour", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting productivity by hour", "user_id", userID, "start_date", startDate, "end_date", endDate)

	var productivity map[string]float64

	// Read operation to get productivity by hour
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		hourlyData, err := qtx.GetProductivityByHour(ctx, store.GetProductivityByHourParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
		})
		if err != nil {
			return fmt.Errorf("failed to get productivity by hour: %w", err)
		}

		productivity = make(map[string]float64)
		for _, hour := range hourlyData {
			// Convert pgtype.Numeric to string for hour key
			hourValue := "0"
			if hour.HourOfDay.Valid {
				// Convert numeric to int64 and then to string
				if val, err := hour.HourOfDay.Value(); err == nil {
					if intVal, ok := val.(int64); ok {
						hourValue = fmt.Sprintf("%d", intVal)
					}
				}
			}
			hourKey := fmt.Sprintf("hour_%s", hourValue)
			productivity[hourKey] = float64(hour.TotalMinutes)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get productivity by hour", "user_id", userID, "start_date", startDate, "end_date", endDate)
		return nil, fmt.Errorf("failed to get productivity by hour: %w", err)
	}

	s.logger.Info("Successfully retrieved productivity by hour", "user_id", userID, "hours_count", len(productivity))
	return productivity, nil
}

// GetTopProjectsByTime retrieves top projects by time spent
func (s *AnalyticsService) GetTopProjectsByTime(ctx context.Context, userID string, startDate, endDate time.Time, limit int32) ([]map[string]any, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetTopProjectsByTime", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	s.logger.Info("Getting top projects by time", "user_id", userID, "start_date", startDate, "end_date", endDate, "limit", limit)

	var projects []map[string]any

	// Read operation to get top projects by time
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		projectData, err := qtx.GetTopProjectsByTime(ctx, store.GetTopProjectsByTimeParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(startDate),
			StartTime_2: timeToPgTimestamptz(endDate),
			Limit:       limit,
		})
		if err != nil {
			return fmt.Errorf("failed to get top projects by time: %w", err)
		}

		projects = make([]map[string]any, len(projectData))
		for i, project := range projectData {
			projects[i] = map[string]any{
				"project_id":    project.ID,
				"project_name":  project.Name,
				"project_color": pgTextToStringRequired(project.Color),
				"total_minutes": project.TotalMinutes,
				"entry_count":   project.EntryCount,
				"percentage":    project.Percentage,
			}
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get top projects by time", "user_id", userID, "start_date", startDate, "end_date", endDate, "limit", limit)
		return nil, fmt.Errorf("failed to get top projects by time: %w", err)
	}

	s.logger.Info("Successfully retrieved top projects by time", "user_id", userID, "projects_count", len(projects), "limit", limit)
	return projects, nil
}

// GetComparisonStats retrieves comparison statistics between two periods
func (s *AnalyticsService) GetComparisonStats(ctx context.Context, userID string, currentStart, currentEnd, previousStart, previousEnd time.Time) (map[string]any, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetComparisonStats", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting comparison stats",
		"user_id", userID,
		"current_start", currentStart,
		"current_end", currentEnd,
		"previous_start", previousStart,
		"previous_end", previousEnd)

	var stats map[string]any

	// Read operation to get comparison stats
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		comparisonData, err := qtx.GetComparisonStats(ctx, store.GetComparisonStatsParams{
			UserID:      userUUID,
			StartTime:   timeToPgTimestamptz(currentStart),
			StartTime_2: timeToPgTimestamptz(currentEnd),
			StartTime_3: timeToPgTimestamptz(previousStart),
			StartTime_4: timeToPgTimestamptz(previousEnd),
		})
		if err != nil {
			return fmt.Errorf("failed to get comparison stats: %w", err)
		}

		// Calculate changes manually
		entriesChange := comparisonData.CurrentEntries - comparisonData.PreviousEntries
		minutesChange := comparisonData.CurrentMinutes - comparisonData.PreviousMinutes

		stats = map[string]any{
			"current_entries":  comparisonData.CurrentEntries,
			"current_minutes":  comparisonData.CurrentMinutes,
			"previous_entries": comparisonData.PreviousEntries,
			"previous_minutes": comparisonData.PreviousMinutes,
			"entries_change":   entriesChange,
			"minutes_change":   minutesChange,
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get comparison stats", "user_id", userID)
		return nil, fmt.Errorf("failed to get comparison stats: %w", err)
	}

	s.logger.Info("Successfully retrieved comparison stats",
		"user_id", userID,
		"current_entries", stats["current_entries"],
		"current_minutes", stats["current_minutes"],
		"entries_change", stats["entries_change"],
		"minutes_change", stats["minutes_change"])

	return stats, nil
}

// RefreshUserActivitySummary refreshes the materialized view for user activity summary
func (s *AnalyticsService) RefreshUserActivitySummary(ctx context.Context) error {
	s.logger.Info("Refreshing user activity summary materialized view")

	// Write operation to refresh materialized view
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.RefreshUserActivitySummary(ctx)
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to refresh user activity summary materialized view")
		return fmt.Errorf("failed to refresh user activity summary: %w", err)
	}

	s.logger.Info("Successfully refreshed user activity summary materialized view")
	return nil
}

// validateDateRange validates that the date range is valid
func (s *AnalyticsService) validateDateRange(startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		s.logger.Warn("Invalid date range provided", "start_date", startDate, "end_date", endDate)
		return fmt.Errorf("end date must be after or equal to start date")
	}

	return nil
}

// validateUserID validates that the user ID is a valid UUID
func (s *AnalyticsService) validateUserID(ctx context.Context, userID string) error {
	if userID == "" {
		s.logger.Warn("Empty user ID provided")
		return fmt.Errorf("invalid user ID: cannot be empty")
	}

	_, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return nil
}
