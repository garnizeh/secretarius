package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReportType_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		reportType ReportType
		want       bool
	}{
		{"valid daily summary", ReportDailySummary, true},
		{"valid weekly summary", ReportWeeklySummary, true},
		{"valid monthly summary", ReportMonthlySummary, true},
		{"valid quarterly summary", ReportQuarterlySummary, true},
		{"valid project analysis", ReportProjectAnalysis, true},
		{"valid productivity trends", ReportProductivityTrends, true},
		{"valid time distribution", ReportTimeDistribution, true},
		{"valid performance review", ReportPerformanceReview, true},
		{"valid goal progress", ReportGoalProgress, true},
		{"valid custom", ReportCustom, true},
		{"invalid type", ReportType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.reportType.IsValid())
		})
	}
}

func TestInsightStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status InsightStatus
		want   bool
	}{
		{"valid active", InsightActive, true},
		{"valid archived", InsightArchived, true},
		{"valid superseded", InsightSuperseded, true},
		{"invalid status", InsightStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestTaskType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		taskType TaskType
		want     bool
	}{
		{"valid generate insight", TaskGenerateInsight, true},
		{"valid send email", TaskSendEmail, true},
		{"valid export data", TaskExportData, true},
		{"valid cleanup data", TaskCleanupData, true},
		{"valid process analytics", TaskProcessAnalytics, true},
		{"valid generate report", TaskGenerateReport, true},
		{"valid backup data", TaskBackupData, true},
		{"valid custom", TaskCustom, true},
		{"invalid type", TaskType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.taskType.IsValid())
		})
	}
}

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{"valid pending", TaskPending, true},
		{"valid processing", TaskProcessing, true},
		{"valid completed", TaskCompleted, true},
		{"valid failed", TaskFailed, true},
		{"valid cancelled", TaskCancelled, true},
		{"valid retrying", TaskRetrying, true},
		{"invalid status", TaskStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestGeneratedInsight_Validate(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	validScore := 0.85

	tests := []struct {
		name    string
		insight GeneratedInsight
		wantErr bool
	}{
		{
			name: "valid insight",
			insight: GeneratedInsight{
				PeriodStart:  baseTime,
				PeriodEnd:    baseTime.Add(24 * time.Hour),
				ReportType:   ReportDailySummary,
				Status:       InsightActive,
				QualityScore: &validScore,
			},
			wantErr: false,
		},
		{
			name: "invalid time range",
			insight: GeneratedInsight{
				PeriodStart: baseTime.Add(24 * time.Hour),
				PeriodEnd:   baseTime,
				ReportType:  ReportDailySummary,
				Status:      InsightActive,
			},
			wantErr: true,
		},
		{
			name: "invalid report type",
			insight: GeneratedInsight{
				PeriodStart: baseTime,
				PeriodEnd:   baseTime.Add(24 * time.Hour),
				ReportType:  ReportType("invalid"),
				Status:      InsightActive,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			insight: GeneratedInsight{
				PeriodStart: baseTime,
				PeriodEnd:   baseTime.Add(24 * time.Hour),
				ReportType:  ReportDailySummary,
				Status:      InsightStatus("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid quality score - too low",
			insight: GeneratedInsight{
				PeriodStart:  baseTime,
				PeriodEnd:    baseTime.Add(24 * time.Hour),
				ReportType:   ReportDailySummary,
				Status:       InsightActive,
				QualityScore: &[]float64{-0.1}[0],
			},
			wantErr: true,
		},
		{
			name: "invalid quality score - too high",
			insight: GeneratedInsight{
				PeriodStart:  baseTime,
				PeriodEnd:    baseTime.Add(24 * time.Hour),
				ReportType:   ReportDailySummary,
				Status:       InsightActive,
				QualityScore: &[]float64{1.1}[0],
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.insight.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: Task{
				TaskType: TaskGenerateInsight,
				Status:   TaskPending,
				Priority: 5,
			},
			wantErr: false,
		},
		{
			name: "invalid task type",
			task: Task{
				TaskType: TaskType("invalid"),
				Status:   TaskPending,
				Priority: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			task: Task{
				TaskType: TaskGenerateInsight,
				Status:   TaskStatus("invalid"),
				Priority: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid priority - too low",
			task: Task{
				TaskType: TaskGenerateInsight,
				Status:   TaskPending,
				Priority: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid priority - too high",
			task: Task{
				TaskType: TaskGenerateInsight,
				Status:   TaskPending,
				Priority: 11,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
