package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReportType represents the type of report/insight being generated
type ReportType string

const (
	ReportDailySummary       ReportType = "daily_summary"
	ReportWeeklySummary      ReportType = "weekly_summary"
	ReportMonthlySummary     ReportType = "monthly_summary"
	ReportQuarterlySummary   ReportType = "quarterly_summary"
	ReportProjectAnalysis    ReportType = "project_analysis"
	ReportProductivityTrends ReportType = "productivity_trends"
	ReportTimeDistribution   ReportType = "time_distribution"
	ReportPerformanceReview  ReportType = "performance_review"
	ReportGoalProgress       ReportType = "goal_progress"
	ReportCustom             ReportType = "custom"
)

// IsValid checks if the ReportType is valid
func (r ReportType) IsValid() bool {
	switch r {
	case ReportDailySummary, ReportWeeklySummary, ReportMonthlySummary, ReportQuarterlySummary,
		ReportProjectAnalysis, ReportProductivityTrends, ReportTimeDistribution,
		ReportPerformanceReview, ReportGoalProgress, ReportCustom:
		return true
	}
	return false
}

// InsightStatus represents the status of a generated insight
type InsightStatus string

const (
	InsightActive     InsightStatus = "active"
	InsightArchived   InsightStatus = "archived"
	InsightSuperseded InsightStatus = "superseded"
)

// IsValid checks if the InsightStatus is valid
func (s InsightStatus) IsValid() bool {
	switch s {
	case InsightActive, InsightArchived, InsightSuperseded:
		return true
	}
	return false
}

// GeneratedInsight represents an AI-generated insight or report
type GeneratedInsight struct {
	ID                   uuid.UUID      `json:"id" db:"id"`
	UserID               uuid.UUID      `json:"user_id" db:"user_id"`
	ReportType           ReportType     `json:"report_type" db:"report_type" validate:"required"`
	PeriodStart          time.Time      `json:"period_start" db:"period_start" validate:"required"`
	PeriodEnd            time.Time      `json:"period_end" db:"period_end" validate:"required"`
	Title                string         `json:"title" db:"title" validate:"required,max=200"`
	Content              string         `json:"content" db:"content" validate:"required"`
	Summary              *string        `json:"summary,omitempty" db:"summary"`
	Metadata             map[string]any `json:"metadata" db:"metadata"`
	GenerationModel      *string        `json:"generation_model,omitempty" db:"generation_model"`
	GenerationDurationMs *int           `json:"generation_duration_ms,omitempty" db:"generation_duration_ms"`
	QualityScore         *float64       `json:"quality_score,omitempty" db:"quality_score" validate:"omitempty,min=0,max=1"`
	Status               InsightStatus  `json:"status" db:"status" validate:"required"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
}

// TaskType represents the type of background task
type TaskType string

const (
	TaskGenerateInsight  TaskType = "generate_insight"
	TaskSendEmail        TaskType = "send_email"
	TaskExportData       TaskType = "export_data"
	TaskCleanupData      TaskType = "cleanup_data"
	TaskProcessAnalytics TaskType = "process_analytics"
	TaskGenerateReport   TaskType = "generate_report"
	TaskBackupData       TaskType = "backup_data"
	TaskCustom           TaskType = "custom"
)

// IsValid checks if the TaskType is valid
func (t TaskType) IsValid() bool {
	switch t {
	case TaskGenerateInsight, TaskSendEmail, TaskExportData, TaskCleanupData,
		TaskProcessAnalytics, TaskGenerateReport, TaskBackupData, TaskCustom:
		return true
	}
	return false
}

// TaskStatus represents the status of a background task
type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskProcessing TaskStatus = "processing"
	TaskCompleted  TaskStatus = "completed"
	TaskFailed     TaskStatus = "failed"
	TaskCancelled  TaskStatus = "cancelled"
	TaskRetrying   TaskStatus = "retrying"
)

// IsValid checks if the TaskStatus is valid
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskPending, TaskProcessing, TaskCompleted, TaskFailed, TaskCancelled, TaskRetrying:
		return true
	}
	return false
}

// Task represents a background task
type Task struct {
	ID                   uuid.UUID      `json:"id" db:"id"`
	TaskType             TaskType       `json:"task_type" db:"task_type" validate:"required"`
	UserID               *uuid.UUID     `json:"user_id,omitempty" db:"user_id"`
	Payload              map[string]any `json:"payload" db:"payload"`
	Status               TaskStatus     `json:"status" db:"status" validate:"required"`
	Priority             int            `json:"priority" db:"priority" validate:"min=1,max=10"`
	MaxRetries           int            `json:"max_retries" db:"max_retries"`
	RetryCount           int            `json:"retry_count" db:"retry_count"`
	ScheduledAt          time.Time      `json:"scheduled_at" db:"scheduled_at"`
	StartedAt            *time.Time     `json:"started_at,omitempty" db:"started_at"`
	CompletedAt          *time.Time     `json:"completed_at,omitempty" db:"completed_at"`
	Result               map[string]any `json:"result,omitempty" db:"result"`
	ErrorMessage         *string        `json:"error_message,omitempty" db:"error_message"`
	ProcessingDurationMs *int           `json:"processing_duration_ms,omitempty" db:"processing_duration_ms"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
}

// Validate validates the generated insight data
func (g *GeneratedInsight) Validate() error {
	if g.PeriodEnd.Before(g.PeriodStart) {
		return errors.New("period_end must be after period_start")
	}
	if !g.ReportType.IsValid() {
		return errors.New("invalid report type")
	}
	if !g.Status.IsValid() {
		return errors.New("invalid status")
	}
	if g.QualityScore != nil && (*g.QualityScore < 0 || *g.QualityScore > 1) {
		return errors.New("quality_score must be between 0 and 1")
	}
	return nil
}

// Validate validates the task data
func (t *Task) Validate() error {
	if !t.TaskType.IsValid() {
		return errors.New("invalid task type")
	}
	if !t.Status.IsValid() {
		return errors.New("invalid task status")
	}
	if t.Priority < 1 || t.Priority > 10 {
		return errors.New("priority must be between 1 and 10")
	}
	return nil
}
