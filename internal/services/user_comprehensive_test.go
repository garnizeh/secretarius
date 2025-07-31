package services

import (
	"fmt"
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestUserService_BusinessLogicScenarios tests business logic edge cases
func TestUserService_BusinessLogicScenarios(t *testing.T) {
	userService := &UserService{}

	t.Run("profile_request_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name    string
			req     *models.UserProfileRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "single_character_names",
				req: &models.UserProfileRequest{
					FirstName: "A",
					LastName:  "B",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "names_with_spaces",
				req: &models.UserProfileRequest{
					FirstName: "John Paul",
					LastName:  "Smith Jones",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "names_with_special_chars",
				req: &models.UserProfileRequest{
					FirstName: "José-María",
					LastName:  "O'Connor",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "unicode_names",
				req: &models.UserProfileRequest{
					FirstName: "张三",
					LastName:  "李四",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "long_timezone",
				req: &models.UserProfileRequest{
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "America/Argentina/Buenos_Aires",
				},
				wantErr: false,
			},
			{
				name: "whitespace_only_first_name",
				req: &models.UserProfileRequest{
					FirstName: "   ",
					LastName:  "Doe",
					Timezone:  "UTC",
				},
				wantErr: true,
				errMsg:  "first name cannot be whitespace only",
			},
			{
				name: "whitespace_only_last_name",
				req: &models.UserProfileRequest{
					FirstName: "John",
					LastName:  "   ",
					Timezone:  "UTC",
				},
				wantErr: true,
				errMsg:  "last name cannot be whitespace only",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.validateProfileRequest(tc.req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("password_change_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name    string
			req     *models.UserPasswordChangeRequest
			wantErr bool
			errMsg  string
		}{
			{
				name: "exactly_8_chars_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpass1",
					NewPassword:     "newpass1",
				},
				wantErr: false,
			},
			{
				name: "exactly_100_chars_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpassword",
					NewPassword:     stringRepeat("a", 100),
				},
				wantErr: false,
			},
			{
				name: "unicode_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpássword",
					NewPassword:     "newpássword123",
				},
				wantErr: false,
			},
			{
				name: "special_chars_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "old!@#$%^&*()",
					NewPassword:     "new!@#$%^&*()",
				},
				wantErr: false,
			},
			{
				name: "same_passwords",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "samepassword123",
					NewPassword:     "samepassword123",
				},
				wantErr: false,
			},
			{
				name: "7_chars_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpass",
					NewPassword:     "newpass",
				},
				wantErr: true,
				errMsg:  "new password must be at least 8 characters long",
			},
			{
				name: "101_chars_password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpassword",
					NewPassword:     stringRepeat("a", 101),
				},
				wantErr: true,
				errMsg:  "new password must be at most 100 characters long",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.validatePasswordChangeRequest(tc.req)
				if tc.wantErr {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.errMsg)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("registration_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name    string
			req     *models.UserRegistration
			wantErr bool
			errMsg  string
		}{
			{
				name: "minimal_valid_data",
				req: &models.UserRegistration{
					Email:     "a@b.co",
					Password:  "password",
					FirstName: "A",
					LastName:  "B",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "long_email",
				req: &models.UserRegistration{
					Email:     "very.long.email.address.with.many.characters@example-domain-with-long-name.com",
					Password:  "password123",
					FirstName: "Long",
					LastName:  "Email",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "unicode_in_email_domain",
				req: &models.UserRegistration{
					Email:     "test@münchen.de",
					Password:  "password123",
					FirstName: "Unicode",
					LastName:  "Domain",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "special_chars_in_email",
				req: &models.UserRegistration{
					Email:     "test+tag@example.com",
					Password:  "password123",
					FirstName: "Special",
					LastName:  "Email",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "whitespace_in_email",
				req: &models.UserRegistration{
					Email:     " test@example.com ",
					Password:  "password123",
					FirstName: "Space",
					LastName:  "Email",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "empty_string_email",
				req: &models.UserRegistration{
					Email:     "",
					Password:  "password123",
					FirstName: "Empty",
					LastName:  "Email",
					Timezone:  "UTC",
				},
				wantErr: true,
				errMsg:  "email is required",
			},
			{
				name: "whitespace_only_email",
				req: &models.UserRegistration{
					Email:     "   ",
					Password:  "password123",
					FirstName: "Whitespace",
					LastName:  "Email",
					Timezone:  "UTC",
				},
				wantErr: true,
				errMsg:  "email cannot be whitespace only",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.validateRegistrationRequest(tc.req)
				if tc.wantErr {
					assert.Error(t, err)
					if err != nil && tc.errMsg != "" {
						assert.Contains(t, err.Error(), tc.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("preferences_serialization_scenarios", func(t *testing.T) {
		testCases := []struct {
			name        string
			preferences map[string]any
			expectError bool
		}{
			{
				name:        "nil_preferences",
				preferences: nil,
				expectError: false,
			},
			{
				name:        "empty_preferences",
				preferences: map[string]any{},
				expectError: false,
			},
			{
				name: "simple_preferences",
				preferences: map[string]any{
					"theme": "dark",
					"lang":  "en",
				},
				expectError: false,
			},
			{
				name: "nested_preferences",
				preferences: map[string]any{
					"ui": map[string]any{
						"theme":   "dark",
						"sidebar": true,
					},
					"notifications": map[string]any{
						"email": false,
						"push":  true,
					},
				},
				expectError: false,
			},
			{
				name: "array_preferences",
				preferences: map[string]any{
					"languages":   []string{"en", "es", "fr"},
					"numbers":     []int{1, 2, 3},
					"mixed_array": []any{"string", 42, true},
				},
				expectError: false,
			},
			{
				name: "complex_nested_structure",
				preferences: map[string]any{
					"dashboard": map[string]any{
						"widgets": []map[string]any{
							{"type": "calendar", "position": 1},
							{"type": "tasks", "position": 2},
						},
						"layout": map[string]any{
							"columns": 3,
							"responsive": map[string]any{
								"mobile":  1,
								"tablet":  2,
								"desktop": 3,
							},
						},
					},
				},
				expectError: false,
			},
			{
				name: "unicode_values",
				preferences: map[string]any{
					"displayName": "José María",
					"locale":      "es-ES",
					"currency":    "€",
				},
				expectError: false,
			},
			{
				name: "numeric_values",
				preferences: map[string]any{
					"version":     1.5,
					"maxItems":    100,
					"timeout":     30.5,
					"enabled":     true,
					"nullable":    nil,
					"largeNumber": 9223372036854775807,
				},
				expectError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test serialization
				bytes, err := userService.preferencesToBytes(tc.preferences)
				if tc.expectError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				assert.NotNil(t, bytes)

				// Test deserialization
				recovered, err := userService.bytesToPreferences(bytes)
				assert.NoError(t, err)
				assert.NotNil(t, recovered)

				// For nil input, we expect empty map output
				if tc.preferences == nil {
					assert.Empty(t, recovered)
					return
				}

				// For empty map, we expect empty map output
				if len(tc.preferences) == 0 {
					assert.Empty(t, recovered)
					return
				}

				// Check that top-level keys match
				assert.Equal(t, len(tc.preferences), len(recovered))
				for key := range tc.preferences {
					assert.Contains(t, recovered, key)
				}
			})
		}
	})

	t.Run("invalid_json_handling", func(t *testing.T) {
		testCases := []struct {
			name    string
			data    []byte
			wantErr bool
		}{
			{
				name:    "valid_empty_object",
				data:    []byte("{}"),
				wantErr: false,
			},
			{
				name:    "valid_json",
				data:    []byte(`{"theme": "dark"}`),
				wantErr: false,
			},
			{
				name:    "invalid_json_missing_quote",
				data:    []byte(`{theme: "dark"}`),
				wantErr: true,
			},
			{
				name:    "invalid_json_trailing_comma",
				data:    []byte(`{"theme": "dark",}`),
				wantErr: true,
			},
			{
				name:    "invalid_json_unmatched_brace",
				data:    []byte(`{"theme": "dark"`),
				wantErr: true,
			},
			{
				name:    "invalid_json_array_as_root",
				data:    []byte(`["theme", "dark"]`),
				wantErr: true,
			},
			{
				name:    "invalid_json_string_as_root",
				data:    []byte(`"just a string"`),
				wantErr: true,
			},
			{
				name:    "completely_invalid_data",
				data:    []byte("this is not json at all"),
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				prefs, err := userService.bytesToPreferences(tc.data)

				// Always expect a non-nil preferences map
				assert.NotNil(t, prefs)

				if tc.wantErr {
					assert.Error(t, err)
					assert.Empty(t, prefs)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// TestUserService_ValidationBoundaryConditions tests validation at boundaries
func TestUserService_ValidationBoundaryConditions(t *testing.T) {
	userService := &UserService{}

	t.Run("password_length_boundaries", func(t *testing.T) {
		testCases := []struct {
			length  int
			wantErr bool
		}{
			{0, true},
			{1, true},
			{7, true},
			{8, false},   // minimum valid
			{50, false},  // middle range
			{100, false}, // maximum valid
			{101, true},  // over maximum
			{200, true},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("password_length_%d", tc.length), func(t *testing.T) {
				password := stringRepeat("a", tc.length)

				// Test registration validation
				regReq := &models.UserRegistration{
					Email:     "test@example.com",
					Password:  password,
					FirstName: "Test",
					LastName:  "User",
					Timezone:  "UTC",
				}

				err := userService.validateRegistrationRequest(regReq)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				// Test password change validation (new password)
				if tc.length > 0 {
					pwChangeReq := &models.UserPasswordChangeRequest{
						CurrentPassword: "oldpassword123",
						NewPassword:     password,
					}

					err = userService.validatePasswordChangeRequest(pwChangeReq)
					if tc.wantErr {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				}
			})
		}
	})

	t.Run("email_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name    string
			email   string
			wantErr bool
		}{
			{
				name:    "shortest_valid_email",
				email:   "a@b.c",
				wantErr: false,
			},
			{
				name:    "email_with_numbers",
				email:   "user123@domain456.com",
				wantErr: false,
			},
			{
				name:    "email_with_dots_in_local_part",
				email:   "first.last@example.com",
				wantErr: false,
			},
			{
				name:    "email_with_plus_in_local_part",
				email:   "user+tag@example.com",
				wantErr: false,
			},
			{
				name:    "email_with_hyphen_in_domain",
				email:   "user@sub-domain.example.com",
				wantErr: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				regReq := &models.UserRegistration{
					Email:     tc.email,
					Password:  "password123",
					FirstName: "Test",
					LastName:  "User",
					Timezone:  "UTC",
				}

				err := userService.validateRegistrationRequest(regReq)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("timezone_edge_cases", func(t *testing.T) {
		testCases := []struct {
			name     string
			timezone string
			wantErr  bool
		}{
			{
				name:     "utc",
				timezone: "UTC",
				wantErr:  false,
			},
			{
				name:     "short_timezone",
				timezone: "EST",
				wantErr:  false,
			},
			{
				name:     "long_timezone",
				timezone: "America/Argentina/Buenos_Aires",
				wantErr:  false,
			},
			{
				name:     "timezone_with_slash",
				timezone: "Europe/London",
				wantErr:  false,
			},
			{
				name:     "timezone_with_underscore",
				timezone: "America/New_York",
				wantErr:  false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				regReq := &models.UserRegistration{
					Email:     "test@example.com",
					Password:  "password123",
					FirstName: "Test",
					LastName:  "User",
					Timezone:  tc.timezone,
				}

				err := userService.validateRegistrationRequest(regReq)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
