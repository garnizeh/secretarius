package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ActivityType represents the type of activity being logged
type ActivityType string

const (
	ActivityDevelopment   ActivityType = "development"
	ActivityMeeting       ActivityType = "meeting"
	ActivityCodeReview    ActivityType = "code_review"
	ActivityDebugging     ActivityType = "debugging"
	ActivityDocumentation ActivityType = "documentation"
	ActivityTesting       ActivityType = "testing"
	ActivityDeployment    ActivityType = "deployment"
	ActivityResearch      ActivityType = "research"
	ActivityPlanning      ActivityType = "planning"
	ActivityLearning      ActivityType = "learning"
	ActivityMaintenance   ActivityType = "maintenance"
	ActivitySupport       ActivityType = "support"
	ActivityOther         ActivityType = "other"
)

// IsValid checks if the ActivityType is valid
func (a ActivityType) IsValid() bool {
	switch a {
	case ActivityDevelopment, ActivityMeeting, ActivityCodeReview, ActivityDebugging,
		ActivityDocumentation, ActivityTesting, ActivityDeployment, ActivityResearch,
		ActivityPlanning, ActivityLearning, ActivityMaintenance, ActivitySupport, ActivityOther:
		return true
	}
	return false
}

// ValueRating represents the perceived value of an activity
type ValueRating string

const (
	ValueLow      ValueRating = "low"
	ValueMedium   ValueRating = "medium"
	ValueHigh     ValueRating = "high"
	ValueCritical ValueRating = "critical"
)

// IsValid checks if the ValueRating is valid
func (v ValueRating) IsValid() bool {
	switch v {
	case ValueLow, ValueMedium, ValueHigh, ValueCritical:
		return true
	}
	return false
}

// ImpactLevel represents the scope of impact for an activity
type ImpactLevel string

const (
	ImpactPersonal   ImpactLevel = "personal"
	ImpactTeam       ImpactLevel = "team"
	ImpactDepartment ImpactLevel = "department"
	ImpactCompany    ImpactLevel = "company"
)

// IsValid checks if the ImpactLevel is valid
func (i ImpactLevel) IsValid() bool {
	switch i {
	case ImpactPersonal, ImpactTeam, ImpactDepartment, ImpactCompany:
		return true
	}
	return false
}

// LogEntry represents a single activity log entry
type LogEntry struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	UserID          uuid.UUID    `json:"user_id" db:"user_id"`
	Title           string       `json:"title" db:"title" validate:"required,max=500"`
	Description     *string      `json:"description,omitempty" db:"description"`
	Type            ActivityType `json:"type" db:"type" validate:"required"`
	ProjectID       *uuid.UUID   `json:"project_id,omitempty" db:"project_id"`
	StartTime       time.Time    `json:"start_time" db:"start_time" validate:"required"`
	EndTime         time.Time    `json:"end_time" db:"end_time" validate:"required"`
	DurationMinutes int          `json:"duration_minutes" db:"duration_minutes"`
	ValueRating     ValueRating  `json:"value_rating" db:"value_rating" validate:"required"`
	ImpactLevel     ImpactLevel  `json:"impact_level" db:"impact_level" validate:"required"`
	Tags            []string     `json:"tags,omitempty"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" db:"updated_at"`
}

// LogEntryRequest represents the data required to create or update a log entry
type LogEntryRequest struct {
	Title       string       `json:"title" validate:"required,max=500"`
	Description *string      `json:"description,omitempty" validate:"omitempty,max=5000"`
	Type        ActivityType `json:"type" validate:"required"`
	ProjectID   *uuid.UUID   `json:"project_id,omitempty"`
	StartTime   time.Time    `json:"start_time" validate:"required"`
	EndTime     time.Time    `json:"end_time" validate:"required"`
	ValueRating ValueRating  `json:"value_rating" validate:"required"`
	ImpactLevel ImpactLevel  `json:"impact_level" validate:"required"`
	Tags        []string     `json:"tags,omitempty" validate:"omitempty,dive,max=100"`
}

// CalculateDuration calculates and sets the duration in minutes for a log entry
func (l *LogEntry) CalculateDuration() {
	l.DurationMinutes = int(l.EndTime.Sub(l.StartTime).Minutes())
}

// Validate validates the log entry data
func (l *LogEntry) Validate() error {
	if l.EndTime.Before(l.StartTime) {
		return errors.New("end_time must be after start_time")
	}
	if !l.Type.IsValid() {
		return errors.New("invalid activity type")
	}
	if !l.ValueRating.IsValid() {
		return errors.New("invalid value rating")
	}
	if !l.ImpactLevel.IsValid() {
		return errors.New("invalid impact level")
	}
	return nil
}

// ProductivityMetrics represents aggregated productivity metrics
type ProductivityMetrics struct {
	TotalActivities     int                  `json:"total_activities"`
	TotalMinutes        int                  `json:"total_minutes"`
	ProjectsWorked      int                  `json:"projects_worked"`
	HighValueActivities int                  `json:"high_value_activities"`
	ActivityBreakdown   map[ActivityType]int `json:"activity_breakdown"`
	ValueDistribution   map[ValueRating]int  `json:"value_distribution"`
	ImpactDistribution  map[ImpactLevel]int  `json:"impact_distribution"`
	DailyAverages       map[string]float64   `json:"daily_averages"`
}

// ActivitySummary represents a summary of activities by type
type ActivitySummary struct {
	Type         ActivityType `json:"type"`
	Count        int          `json:"count"`
	TotalMinutes int          `json:"total_minutes"`
	AvgMinutes   float64      `json:"avg_minutes"`
	Week         time.Time    `json:"week"`
}
