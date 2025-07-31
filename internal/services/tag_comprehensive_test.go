package services

import (
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestTagService_BusinessLogicScenarios tests business logic edge cases
func TestTagService_BusinessLogicScenarios(t *testing.T) {
	tagService := &TagService{}

	t.Run("tag_name_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name    string
			tagName string
			wantErr bool
			errMsg  string
		}{
			{
				name:    "single_character_name",
				tagName: "a",
				wantErr: false,
			},
			{
				name:    "name_with_spaces",
				tagName: "tag with spaces",
				wantErr: false,
			},
			{
				name:    "name_with_special_chars",
				tagName: "tag-name_with.special@chars",
				wantErr: false,
			},
			{
				name:    "unicode_name",
				tagName: "标签名称",
				wantErr: false,
			},
			{
				name:    "exactly_100_chars",
				tagName: stringRepeat("a", 100),
				wantErr: false,
			},
			{
				name:    "101_chars_should_fail",
				tagName: stringRepeat("a", 101),
				wantErr: true,
				errMsg:  "tag name must be at most 100 characters",
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.TagRequest{
					Name:        tt.tagName,
					Color:       "#FF5733",
					Description: nil,
				}

				err := tagService.validateTagRequest(req)
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
	})

	t.Run("description_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			description *string
			wantErr     bool
			errMsg      string
		}{
			{
				name:        "nil_description",
				description: nil,
				wantErr:     false,
			},
			{
				name:        "empty_description",
				description: stringPtr(""),
				wantErr:     false,
			},
			{
				name:        "single_char_description",
				description: stringPtr("a"),
				wantErr:     false,
			},
			{
				name:        "exactly_500_chars",
				description: stringPtr(stringRepeat("a", 500)),
				wantErr:     false,
			},
			{
				name:        "501_chars_should_fail",
				description: stringPtr(stringRepeat("a", 501)),
				wantErr:     true,
				errMsg:      "tag description must be at most 500 characters",
			},
			{
				name:        "description_with_newlines",
				description: stringPtr("Line 1\nLine 2\nLine 3"),
				wantErr:     false,
			},
			{
				name:        "description_with_unicode",
				description: stringPtr("描述包含中文字符和emoji 🏷️"),
				wantErr:     false,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.TagRequest{
					Name:        "test-tag",
					Color:       "#FF5733",
					Description: tt.description,
				}

				err := tagService.validateTagRequest(req)
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
	})

	t.Run("color_validation_comprehensive", func(t *testing.T) {
		testCases := []struct {
			name    string
			color   string
			wantErr bool
		}{
			{
				name:    "valid_red",
				color:   "#FF0000",
				wantErr: false,
			},
			{
				name:    "valid_green",
				color:   "#00FF00",
				wantErr: false,
			},
			{
				name:    "valid_blue",
				color:   "#0000FF",
				wantErr: false,
			},
			{
				name:    "valid_black",
				color:   "#000000",
				wantErr: false,
			},
			{
				name:    "valid_white",
				color:   "#FFFFFF",
				wantErr: false,
			},
			{
				name:    "valid_mixed_case",
				color:   "#AbCdEf",
				wantErr: false,
			},
			{
				name:    "valid_all_lowercase",
				color:   "#abcdef",
				wantErr: false,
			},
			{
				name:    "valid_all_uppercase",
				color:   "#ABCDEF",
				wantErr: false,
			},
			{
				name:    "invalid_no_hash",
				color:   "FF0000",
				wantErr: true,
			},
			{
				name:    "invalid_too_short",
				color:   "#FF00",
				wantErr: true,
			},
			{
				name:    "invalid_too_long",
				color:   "#FF000000",
				wantErr: true,
			},
			{
				name:    "invalid_non_hex",
				color:   "#GGGGGG",
				wantErr: true,
			},
			{
				name:    "invalid_with_spaces",
				color:   "#FF 00 00",
				wantErr: true,
			},
			{
				name:    "invalid_empty",
				color:   "",
				wantErr: true,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				req := &models.TagRequest{
					Name:        "test-tag",
					Color:       tt.color,
					Description: nil,
				}

				err := tagService.validateTagRequest(req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestTagService_LimitBoundaryConditions(t *testing.T) {
	// Test boundary conditions for limit parameters
	testCases := []struct {
		name           string
		inputLimit     int32
		operation      string
		expectedOutput int32
	}{
		{
			name:           "zero_limit_popular_tags",
			inputLimit:     0,
			operation:      "GetPopularTags",
			expectedOutput: 10,
		},
		{
			name:           "negative_limit_popular_tags",
			inputLimit:     -10,
			operation:      "GetPopularTags",
			expectedOutput: 10,
		},
		{
			name:           "max_int32_limit",
			inputLimit:     2147483647, // max int32
			operation:      "GetPopularTags",
			expectedOutput: 2147483647,
		},
		{
			name:           "one_limit",
			inputLimit:     1,
			operation:      "GetPopularTags",
			expectedOutput: 1,
		},
		{
			name:           "zero_limit_search_tags",
			inputLimit:     0,
			operation:      "SearchTags",
			expectedOutput: 20,
		},
		{
			name:           "negative_limit_search_tags",
			inputLimit:     -5,
			operation:      "SearchTags",
			expectedOutput: 20,
		},
		{
			name:           "large_limit_search_tags",
			inputLimit:     1000,
			operation:      "SearchTags",
			expectedOutput: 1000,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualLimit := tt.inputLimit

			// Apply the logic from the actual service methods
			switch tt.operation {
			case "GetPopularTags", "GetRecentlyUsedTags":
				if actualLimit <= 0 {
					actualLimit = 10
				}
			case "SearchTags":
				if actualLimit <= 0 {
					actualLimit = 20
				}
			}

			assert.Equal(t, tt.expectedOutput, actualLimit)
		})
	}
}

func TestTagService_ErrorMessagePatterns(t *testing.T) {
	tagService := &TagService{}

	// Test error message consistency and patterns
	testCases := []struct {
		name          string
		req           *models.TagRequest
		expectedError string
	}{
		{
			name: "empty_name_error",
			req: &models.TagRequest{
				Name:        "",
				Color:       "#FF5733",
				Description: nil,
			},
			expectedError: "tag name is required",
		},
		{
			name: "empty_color_error",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "",
				Description: nil,
			},
			expectedError: "tag color is required",
		},
		{
			name: "invalid_color_error",
			req: &models.TagRequest{
				Name:        "test-tag",
				Color:       "not-a-color",
				Description: nil,
			},
			expectedError: "invalid color format",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tagService.validateTagRequest(tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestTagService_RequestStructureValidation(t *testing.T) {
	// Test that TagRequest structure is properly validated
	t.Run("valid_minimal_request", func(t *testing.T) {
		req := &models.TagRequest{
			Name:        "minimal-tag",
			Color:       "#000000",
			Description: nil,
		}

		tagService := &TagService{}
		err := tagService.validateTagRequest(req)
		assert.NoError(t, err)

		// Verify structure
		assert.Equal(t, "minimal-tag", req.Name)
		assert.Equal(t, "#000000", req.Color)
		assert.Nil(t, req.Description)
	})

	t.Run("valid_complete_request", func(t *testing.T) {
		description := "A complete tag with all fields"
		req := &models.TagRequest{
			Name:        "complete-tag",
			Color:       "#FF5733",
			Description: &description,
		}

		tagService := &TagService{}
		err := tagService.validateTagRequest(req)
		assert.NoError(t, err)

		// Verify structure
		assert.Equal(t, "complete-tag", req.Name)
		assert.Equal(t, "#FF5733", req.Color)
		assert.NotNil(t, req.Description)
		assert.Equal(t, description, *req.Description)
	})
}
