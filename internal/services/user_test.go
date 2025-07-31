package services

import (
	"testing"

	"github.com/garnizeh/englog/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestUserService_ValidationMethods tests validation functions
func TestUserService_ValidationMethods(t *testing.T) {
	userService := &UserService{}

	t.Run("validateProfileRequest", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *models.UserProfileRequest
			wantErr bool
		}{
			{
				name: "valid request",
				req: &models.UserProfileRequest{
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "UTC",
					Preferences: map[string]any{
						"theme": "dark",
					},
				},
				wantErr: false,
			},
			{
				name: "missing first name",
				req: &models.UserProfileRequest{
					FirstName: "",
					LastName:  "Doe",
					Timezone:  "UTC",
				},
				wantErr: true,
			},
			{
				name: "missing last name",
				req: &models.UserProfileRequest{
					FirstName: "John",
					LastName:  "",
					Timezone:  "UTC",
				},
				wantErr: true,
			},
			{
				name: "missing timezone",
				req: &models.UserProfileRequest{
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := userService.validateProfileRequest(tt.req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("validatePasswordChangeRequest", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *models.UserPasswordChangeRequest
			wantErr bool
		}{
			{
				name: "valid request",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpassword123",
					NewPassword:     "newpassword123",
				},
				wantErr: false,
			},
			{
				name: "missing current password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "",
					NewPassword:     "newpassword123",
				},
				wantErr: true,
			},
			{
				name: "missing new password",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpassword123",
					NewPassword:     "",
				},
				wantErr: true,
			},
			{
				name: "new password too short",
				req: &models.UserPasswordChangeRequest{
					CurrentPassword: "oldpassword123",
					NewPassword:     "short",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := userService.validatePasswordChangeRequest(tt.req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("validateRegistrationRequest", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *models.UserRegistration
			wantErr bool
		}{
			{
				name: "valid request",
				req: &models.UserRegistration{
					Email:     "john.doe@example.com",
					Password:  "password123",
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "UTC",
				},
				wantErr: false,
			},
			{
				name: "missing email",
				req: &models.UserRegistration{
					Email:     "",
					Password:  "password123",
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "UTC",
				},
				wantErr: true,
			},
			{
				name: "password too short",
				req: &models.UserRegistration{
					Email:     "john.doe@example.com",
					Password:  "short",
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "UTC",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := userService.validateRegistrationRequest(tt.req)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestUserService_PreferencesConversion(t *testing.T) {
	userService := &UserService{}

	t.Run("ConvertToJSON", func(t *testing.T) {
		preferences := map[string]any{
			"theme":    "dark",
			"language": "en",
			"notifications": map[string]any{
				"email": true,
				"push":  false,
			},
		}

		bytes, err := userService.preferencesToBytes(preferences)
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
		assert.Contains(t, string(bytes), "theme")
		assert.Contains(t, string(bytes), "dark")
	})

	t.Run("ConvertFromJSON", func(t *testing.T) {
		jsonBytes := []byte(`{"theme":"light","language":"pt"}`)

		preferences, err := userService.bytesToPreferences(jsonBytes)
		assert.NoError(t, err)
		assert.Equal(t, "light", preferences["theme"])
		assert.Equal(t, "pt", preferences["language"])
	})

	t.Run("ConvertToBytes", func(t *testing.T) {
		preferences := map[string]any{
			"theme":         "dark",
			"language":      "en",
			"notifications": true,
			"timezone":      "UTC",
		}

		bytes, err := userService.preferencesToBytes(preferences)
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("ConvertFromBytes", func(t *testing.T) {
		jsonBytes := []byte(`{"theme":"dark","language":"en","notifications":true,"timezone":"UTC"}`)

		preferences, err := userService.bytesToPreferences(jsonBytes)
		assert.NoError(t, err)

		assert.Equal(t, "dark", preferences["theme"])
		assert.Equal(t, "en", preferences["language"])
		assert.Equal(t, true, preferences["notifications"])
		assert.Equal(t, "UTC", preferences["timezone"])
	})
}

func TestUserService_PreferencesEdgeCases(t *testing.T) {
	userService := &UserService{}

	t.Run("EmptyPreferences", func(t *testing.T) {
		emptyPrefs := map[string]any{}
		bytes, err := userService.preferencesToBytes(emptyPrefs)
		assert.NoError(t, err)
		assert.Equal(t, []byte("{}"), bytes)
	})

	t.Run("NilPreferences", func(t *testing.T) {
		bytes, err := userService.preferencesToBytes(nil)
		assert.NoError(t, err)
		assert.Equal(t, []byte("{}"), bytes)
	})

	t.Run("EmptyJSONBytes", func(t *testing.T) {
		prefs, err := userService.bytesToPreferences([]byte{})
		assert.NoError(t, err)
		assert.Empty(t, prefs)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSON := []byte("invalid json")
		prefs, err := userService.bytesToPreferences(invalidJSON)
		assert.Error(t, err)
		assert.Empty(t, prefs)
	})
}
