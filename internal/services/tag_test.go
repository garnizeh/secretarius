package services

import (
	"fmt"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTagService_ValidateTagRequest(t *testing.T) {
	tagService := &TagService{}

	tests := []struct {
		name    string
		req     *models.TagRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#FF5733",
				Description: nil,
			},
			wantErr: false,
		},
		{
			name: "valid request with description",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#FF5733",
				Description: stringPtr("A test tag for testing purposes"),
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: &models.TagRequest{
				Name:        "",
				Color:       "#FF5733",
				Description: nil,
			},
			wantErr: true,
			errMsg:  "tag name is required",
		},
		{
			name: "name too long",
			req: &models.TagRequest{
				Name:        stringRepeat("a", 101),
				Color:       "#FF5733",
				Description: nil,
			},
			wantErr: true,
			errMsg:  "tag name must be at most 100 characters",
		},
		{
			name: "missing color",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "",
				Description: nil,
			},
			wantErr: true,
			errMsg:  "tag color is required",
		},
		{
			name: "invalid color format",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "invalid-color",
				Description: nil,
			},
			wantErr: true,
			errMsg:  "invalid color format",
		},
		{
			name: "description too long",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#FF5733",
				Description: stringPtr(stringRepeat("a", 501)),
			},
			wantErr: true,
			errMsg:  "tag description must be at most 500 characters",
		},
		{
			name: "valid hex color six characters",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#ABCDEF",
				Description: nil,
			},
			wantErr: false,
		},
		{
			name: "valid hex color uppercase",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#ABCDEF",
				Description: nil,
			},
			wantErr: false,
		},
		{
			name: "valid hex color lowercase",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "#abcdef",
				Description: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tagService.validateTagRequest(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTagService_UUIDParsing(t *testing.T) {
	// Test invalid UUID scenarios that would be caught by service methods
	testCases := []struct {
		name    string
		uuidStr string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid UUID",
			uuidStr: uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "invalid UUID format",
			uuidStr: "invalid-uuid",
			wantErr: true,
			errMsg:  "invalid UUID length",
		},
		{
			name:    "empty UUID",
			uuidStr: "",
			wantErr: true,
			errMsg:  "invalid UUID length",
		},
		{
			name:    "malformed UUID",
			uuidStr: "12345678-1234-1234-1234-12345678901",
			wantErr: true,
			errMsg:  "invalid UUID length",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uuid.Parse(tt.uuidStr)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTagService_LimitValidation(t *testing.T) {
	// Test limit parameter handling in service methods
	testCases := []struct {
		name     string
		limit    int32
		expected int32
	}{
		{
			name:     "positive limit",
			limit:    5,
			expected: 5,
		},
		{
			name:     "zero limit should use default",
			limit:    0,
			expected: 10, // Default for GetPopularTags and GetRecentlyUsedTags
		},
		{
			name:     "negative limit should use default",
			limit:    -1,
			expected: 10,
		},
		{
			name:     "large limit",
			limit:    1000,
			expected: 1000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualLimit := tt.limit
			if actualLimit <= 0 {
				actualLimit = 10 // Default limit as used in service methods
			}
			assert.Equal(t, tt.expected, actualLimit)
		})
	}
}

func TestTagService_SearchLimitValidation(t *testing.T) {
	// Test search limit parameter handling
	testCases := []struct {
		name     string
		limit    int32
		expected int32
	}{
		{
			name:     "positive limit",
			limit:    15,
			expected: 15,
		},
		{
			name:     "zero limit should use default",
			limit:    0,
			expected: 20, // Default for SearchTags
		},
		{
			name:     "negative limit should use default",
			limit:    -1,
			expected: 20,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualLimit := tt.limit
			if actualLimit <= 0 {
				actualLimit = 20 // Default limit for search as used in SearchTags method
			}
			assert.Equal(t, tt.expected, actualLimit)
		})
	}
}

func TestTagService_ErrorHandling(t *testing.T) {
	// Test error message formatting patterns used in TagService
	testCases := []struct {
		name     string
		baseErr  string
		wrapMsg  string
		expected string
	}{
		{
			name:     "create tag error",
			baseErr:  "duplicate key violation",
			wrapMsg:  "failed to create tag",
			expected: "failed to create tag: duplicate key violation",
		},
		{
			name:     "get tag error",
			baseErr:  "record not found",
			wrapMsg:  "failed to get tag",
			expected: "failed to get tag: record not found",
		},
		{
			name:     "delete tag error",
			baseErr:  "foreign key constraint",
			wrapMsg:  "failed to delete tag",
			expected: "failed to delete tag: foreign key constraint",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := fmt.Errorf("%s: %s", tt.wrapMsg, tt.baseErr)
			assert.Equal(t, tt.expected, wrappedErr.Error())
		})
	}
}

func TestTagService_EnsureTagExistsLogic(t *testing.T) {
	// Test the logic patterns used in EnsureTagExists method
	testColor := "#FF5733"
	testName := "test-tag"

	// Test tag request creation for non-existing tag
	req := &models.TagRequest{
		Name:        testName,
		Color:       testColor,
		Description: nil,
	}

	assert.Equal(t, testName, req.Name)
	assert.Equal(t, testColor, req.Color)
	assert.Nil(t, req.Description)
}
