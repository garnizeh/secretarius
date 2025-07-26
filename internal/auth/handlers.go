package auth

import (
	"net/http"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper function to convert store.User to models.UserProfile
func convertToUserProfile(user store.User) models.UserProfile {
	profile := models.UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	if user.Timezone.Valid {
		profile.Timezone = user.Timezone.String
	}

	// For preferences, since it's []byte, we'll leave it as nil for now
	// In a real implementation, you'd want to marshal/unmarshal JSON
	profile.Preferences = nil

	if user.LastLoginAt.Valid {
		profile.LastLoginAt = &user.LastLoginAt.Time
	}

	profile.CreatedAt = user.CreatedAt.Time

	return profile
}

func (a *AuthService) RegisterHandler(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	_, err := a.queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := a.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password",
		})
		return
	}

	// Create user
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	var preferences []byte // Initialize empty preferences

	user, err := a.queries.CreateUser(c.Request.Context(), store.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Timezone:     pgtype.Text{String: timezone, Valid: true},
		Preferences:  preferences,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Create tokens for immediate login
	accessToken, err := a.CreateAccessToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access token",
		})
		return
	}

	refreshToken, err := a.CreateRefreshToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": convertToUserProfile(user),
		"tokens": models.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(a.accessTokenTTL.Seconds()),
			TokenType:    "Bearer",
		},
	})
}

func (a *AuthService) LoginHandler(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user by email
	user, err := a.queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Check password
	if !a.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Update last login
	err = a.queries.UpdateUserLastLogin(c.Request.Context(), user.ID)
	if err != nil {
		// Log error but don't fail the login
	}

	// Create tokens
	accessToken, err := a.CreateAccessToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access token",
		})
		return
	}

	refreshToken, err := a.CreateRefreshToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": convertToUserProfile(user),
		"tokens": models.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(a.accessTokenTTL.Seconds()),
			TokenType:    "Bearer",
		},
	})
}

func (a *AuthService) RefreshHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	newAccessToken, newRefreshToken, err := a.RotateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid refresh token",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": models.AuthTokens{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
			ExpiresIn:    int(a.accessTokenTTL.Seconds()),
			TokenType:    "Bearer",
		},
	})
}

func (a *AuthService) LogoutHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	claims, err := a.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	if claims.TokenType != models.TokenRefresh {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token type",
		})
		return
	}

	// Denylist the refresh token
	if claims.JTI != "" {
		err = a.DenylistRefreshToken(claims.JTI, claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to logout",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (a *AuthService) MeHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := a.queries.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": convertToUserProfile(user),
	})
}
