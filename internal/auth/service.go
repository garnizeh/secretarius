package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db              *pgxpool.Pool
	queries         *store.Queries
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	UserID    string           `json:"sub"`
	TokenType models.TokenType `json:"type"`
	JTI       string           `json:"jti,omitempty"` // JWT ID for refresh tokens
	jwt.RegisteredClaims
}

func NewAuthService(db *pgxpool.Pool, secretKey string) *AuthService {
	return &AuthService{
		db:              db,
		queries:         store.New(db),
		secretKey:       []byte(secretKey),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 30 * 24 * time.Hour, // 30 days
	}
}

func (a *AuthService) CreateAccessToken(userID string) (string, error) {
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
	return token.SignedString(a.secretKey)
}

func (a *AuthService) CreateRefreshToken(userID string) (string, error) {
	jti, err := generateJTI()
	if err != nil {
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
	return token.SignedString(a.secretKey)
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Check if refresh token is denylisted
	if claims.TokenType == models.TokenRefresh && claims.JTI != "" {
		denylisted, err := a.isRefreshTokenDenylisted(claims.JTI)
		if err != nil {
			return nil, err
		}
		if denylisted {
			return nil, ErrTokenDenylisted
		}
	}

	return claims, nil
}

func (a *AuthService) RotateRefreshToken(oldRefreshToken string) (string, string, error) {
	claims, err := a.ValidateToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.TokenType != models.TokenRefresh {
		return "", "", ErrInvalidTokenType
	}

	// Denylist the old refresh token
	if claims.JTI != "" {
		err = a.DenylistRefreshToken(claims.JTI, claims.UserID)
		if err != nil {
			return "", "", err
		}
	}

	// Generate new tokens
	newAccessToken, err := a.CreateAccessToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := a.CreateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthService) DenylistRefreshToken(jti, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(a.refreshTokenTTL)
	return a.queries.CreateRefreshTokenDenylist(context.Background(), store.CreateRefreshTokenDenylistParams{
		Jti:       jti,
		UserID:    userUUID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
		Column4:   "logout",
	})
}

func (a *AuthService) isRefreshTokenDenylisted(jti string) (bool, error) {
	return a.queries.IsRefreshTokenDenylisted(context.Background(), jti)
}

func (a *AuthService) CleanupExpiredTokens() error {
	return a.queries.CleanupExpiredDenylistedTokens(context.Background())
}

func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
