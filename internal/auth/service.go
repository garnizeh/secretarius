package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// Authentication errors
var (
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenDenylisted      = errors.New("token is denylisted")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserAlreadyExists    = errors.New("user already exists")
)

type AuthService struct {
	db              *database.DB
	logger          *logging.Logger
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	bcryptCost      int
}

type Claims struct {
	UserID    string           `json:"sub"`
	TokenType models.TokenType `json:"type"`
	JTI       string           `json:"jti,omitempty"` // JWT ID for refresh tokens
	jwt.RegisteredClaims
}

// TODO(rodrigo): use key ring or vault for secret key management in production
func NewAuthService(db *database.DB, logger *logging.Logger, secretKey string) *AuthService {
	return &AuthService{
		db:              db,
		logger:          logger.WithComponent("auth_service"),
		secretKey:       []byte(secretKey),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 30 * 24 * time.Hour, // 30 days
		bcryptCost:      bcrypt.DefaultCost,
	}
}

// NewAuthServiceForTest creates an AuthService with lower bcrypt cost for faster testing
func NewAuthServiceForTest(db *database.DB, logger *logging.Logger, secretKey string) *AuthService {
	return &AuthService{
		db:              db,
		logger:          logger.WithComponent("auth_service"),
		secretKey:       []byte(secretKey),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 30 * 24 * time.Hour, // 30 days
		bcryptCost:      4,                   // Much faster for testing
	}
}

func (a *AuthService) CreateAccessToken(ctx context.Context, userID string) (string, error) {
	a.logger.WithContext(ctx).Info("Creating access token", "user_id", userID)

	claims := &Claims{
		UserID:    userID,
		TokenType: models.TokenAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "englog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to sign access token", "user_id", userID)
		return "", err
	}

	a.logger.WithContext(ctx).Info("Access token created successfully", "user_id", userID, "expires_in_minutes", int(a.accessTokenTTL.Minutes()))
	return tokenString, nil
}

func (a *AuthService) CreateRefreshToken(ctx context.Context, userID string) (string, error) {
	a.logger.WithContext(ctx).Info("Creating refresh token", "user_id", userID)

	jti, err := generateJTI()
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to generate JTI for refresh token", "user_id", userID)
		return "", err
	}

	claims := &Claims{
		UserID:    userID,
		TokenType: models.TokenRefresh,
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "englog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to sign refresh token", "user_id", userID, "jti", jti)
		return "", err
	}

	a.logger.WithContext(ctx).Info("Refresh token created successfully", "user_id", userID, "jti", jti, "expires_in_hours", int(a.refreshTokenTTL.Hours()))
	return tokenString, nil
}

func (a *AuthService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	a.logger.WithContext(ctx).Info("Validating token")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return a.secretKey, nil
	})

	if err != nil {
		a.logger.Warn("Token parsing failed", "error", err.Error())
		return nil, err
	}

	if !token.Valid {
		a.logger.Warn("Invalid token provided")
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		a.logger.Warn("Failed to extract claims from token")
		return nil, ErrInvalidToken
	}

	// Check if refresh token is denylisted
	if claims.TokenType == models.TokenRefresh && claims.JTI != "" {
		denylisted, err := a.isRefreshTokenDenylisted(ctx, claims.JTI)
		if err != nil {
			a.logger.LogError(ctx, err, "Failed to check token denylist", "jti", claims.JTI, "user_id", claims.UserID)
			return nil, err
		}
		if denylisted {
			a.logger.Warn("Attempted use of denylisted refresh token", "jti", claims.JTI, "user_id", claims.UserID)
			return nil, ErrTokenDenylisted
		}
	}

	a.logger.WithContext(ctx).Info("Token validated successfully", "user_id", claims.UserID, "token_type", claims.TokenType, "jti", claims.JTI)
	return claims, nil
}

func (a *AuthService) RotateRefreshToken(ctx context.Context, oldRefreshToken string) (string, string, error) {
	a.logger.Info("Rotating refresh token")

	claims, err := a.ValidateToken(ctx, oldRefreshToken)
	if err != nil {
		a.logger.Warn("Failed to validate refresh token for rotation", "error", err.Error())
		return "", "", err
	}

	if claims.TokenType != models.TokenRefresh {
		a.logger.Warn("Invalid token type for rotation", "token_type", claims.TokenType, "user_id", claims.UserID)
		return "", "", ErrInvalidTokenType
	}

	a.logger.Info("Rotating tokens for user", "user_id", claims.UserID, "old_jti", claims.JTI)

	// Denylist the old refresh token
	if claims.JTI != "" {
		err = a.DenylistRefreshToken(ctx, claims.JTI, claims.UserID)
		if err != nil {
			a.logger.LogError(ctx, err, "Failed to denylist old refresh token", "jti", claims.JTI, "user_id", claims.UserID)
			return "", "", err
		}
	}

	// Generate new tokens
	newAccessToken, err := a.CreateAccessToken(ctx, claims.UserID)
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to create new access token during rotation", "user_id", claims.UserID)
		return "", "", err
	}

	newRefreshToken, err := a.CreateRefreshToken(ctx, claims.UserID)
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to create new refresh token during rotation", "user_id", claims.UserID)
		return "", "", err
	}

	a.logger.Info("Token rotation completed successfully", "user_id", claims.UserID, "old_jti", claims.JTI)
	return newAccessToken, newRefreshToken, nil
}

func (a *AuthService) HashPassword(ctx context.Context, password string) (string, error) {
	a.logger.WithContext(ctx).Info("Hashing password")

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		a.logger.LogError(ctx, err, "Failed to hash password")
		return "", err
	}

	a.logger.WithContext(ctx).Info("Password hashed successfully")
	return string(bytes), nil
}

func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	isValid := err == nil

	if isValid {
		a.logger.Info("Password verification successful")
	} else {
		a.logger.Warn("Password verification failed")
	}

	return isValid
}

func (a *AuthService) DenylistRefreshToken(ctx context.Context, jti, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		a.logger.LogError(ctx, err, "Invalid user ID format for token denylist", "user_id", userID, "jti", jti)
		return err
	}

	a.logger.Info("Denylisting refresh token", "user_id", userID, "jti", jti)

	expiresAt := time.Now().Add(a.refreshTokenTTL)
	err = a.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.CreateRefreshTokenDenylist(context.Background(), store.CreateRefreshTokenDenylistParams{
			Jti:       jti,
			UserID:    userUUID,
			ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
			Column4:   "logout",
		})
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			a.logger.Warn("Token already denylisted", "user_id", userID, "jti", jti)
		} else {
			a.logger.LogError(ctx, err, "Failed to denylist refresh token", "user_id", userID, "jti", jti)
		}
		return err
	}

	a.logger.Info("Refresh token denylisted successfully", "user_id", userID, "jti", jti, "expires_at", expiresAt)
	return nil
}

func (a *AuthService) isRefreshTokenDenylisted(ctx context.Context, jti string) (bool, error) {
	a.logger.Info("Checking if refresh token is denylisted", "jti", jti)

	var isDenied bool
	if err := a.db.Read(ctx, func(qtx *store.Queries) error {
		var err error
		isDenied, err = qtx.IsRefreshTokenDenylisted(ctx, jti)
		if err != nil {
			return fmt.Errorf("failed to check denylist: %w", err)
		}
		return nil
	}); err != nil {
		a.logger.LogError(ctx, err, "Failed to check token denylist", "jti", jti)
		return false, err
	}

	if isDenied {
		a.logger.Info("Token is denylisted", "jti", jti)
	} else {
		a.logger.Info("Token is not denylisted", "jti", jti)
	}
	return isDenied, nil
}

func (a *AuthService) CleanupExpiredTokens(ctx context.Context) error {
	a.logger.Info("Starting cleanup of expired tokens")

	err := a.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.CleanupExpiredDenylistedTokens(ctx)
	})

	if err != nil {
		a.logger.LogError(ctx, err, "Failed to cleanup expired tokens")
		return err
	}

	a.logger.Info("Expired tokens cleanup completed successfully")
	return nil
}

func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
