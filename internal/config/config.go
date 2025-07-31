package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the unified application configuration
type Config struct {
	Environment string
	Port        int
	Host        string

	DB        DBConfig
	Auth      AuthConfig
	Server    ServerConfig
	RateLimit RateLimitConfig
	Security  SecurityConfig
	Logging   LoggingConfig
	Redis     RedisConfig

	// gRPC configuration for worker communication
	GRPC   GRPCConfig
	Worker WorkerConfig
}

// DBConfig holds database connection configuration
type DBConfig struct {
	User          string
	Password      string
	HostReadWrite string
	HostReadOnly  string
	Host          string
	Name          string
	Schema        string
}

// AuthConfig holds JWT and authentication configuration
type AuthConfig struct {
	JWTSecretKey    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CleanupInterval time.Duration
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration // Individual request timeout for middleware
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	Enabled           bool
	RedisEnabled      bool
	WindowSize        time.Duration
}

// SecurityConfig holds security and CORS configuration
type SecurityConfig struct {
	CORSAllowedOrigins   []string
	CORSAllowedMethods   []string
	CORSAllowedHeaders   []string
	CORSAllowCredentials bool
	SecurityHeaders      SecurityHeadersConfig
}

// SecurityHeadersConfig holds security headers configuration
type SecurityHeadersConfig struct {
	ContentTypeOptions      string
	FrameOptions            string
	XSSProtection           string
	StrictTransportSecurity string
	ContentSecurityPolicy   string
	ReferrerPolicy          string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      slog.Level
	Format     string // "json" or "text"
	AddSource  bool
	TimeFormat string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// GRPCConfig holds gRPC configuration for worker communication
type GRPCConfig struct {
	ServerPort    int
	WorkerAddress string
	TLSCertFile   string
	TLSKeyFile    string
	TLSEnabled    bool
	// Client configuration for worker
	APIServerAddress string
	ServerName       string
}

// WorkerConfig holds worker-specific configuration
type WorkerConfig struct {
	HealthPort         int    // HTTP port for health check endpoints
	WorkerID           string // Unique worker identifier
	WorkerName         string // Human-readable worker name
	Version            string // Worker version
	OllamaURL          string // Ollama service URL
	MaxConcurrentTasks int    // Maximum concurrent tasks
}

// Load creates and returns a new Config instance with values from environment variables
func Load() *Config {
	cfg := &Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getIntEnv("APP_PORT", 8080),
		Host:        getEnv("APP_HOST", "localhost"),

		DB: DBConfig{
			User:          getEnv("DB_USER", "englog"),
			Password:      getEnv("DB_PASSWORD", "password"),
			HostReadWrite: getEnv("DB_HOST_READ_WRITE", "postgres-dev:5432"),
			HostReadOnly:  getEnv("DB_HOST_READ_ONLY", "postgres-dev:5432"),
			Name:          getEnv("DB_NAME", "englog"),
			Schema:        getEnv("DB_SCHEMA", "englog"),
		},

		Auth: AuthConfig{
			JWTSecretKey:    getEnv("JWT_SECRET_KEY", "your-secret-key-here-change-in-production"),
			AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TOKEN_TTL", 30*24*time.Hour),
			CleanupInterval: getDurationEnv("JWT_DENYLIST_CLEANUP_INTERVAL", 24*time.Hour),
		},

		Server: ServerConfig{
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			RequestTimeout:  getDurationEnv("SERVER_REQUEST_TIMEOUT", 30*time.Second),
		},

		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			BurstSize:         getIntEnv("RATE_LIMIT_BURST", 20),
			Enabled:           getBoolEnv("RATE_LIMIT_ENABLED", true),
			RedisEnabled:      getBoolEnv("RATE_LIMIT_REDIS_ENABLED", true),
			WindowSize:        getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
		},

		Security: SecurityConfig{
			CORSAllowedOrigins:   getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			CORSAllowedMethods:   getSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			CORSAllowedHeaders:   getSliceEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
			CORSAllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
			// Default security headers
			SecurityHeaders: SecurityHeadersConfig{
				ContentTypeOptions:      "nosniff",
				FrameOptions:            "DENY",
				XSSProtection:           "1; mode=block",
				StrictTransportSecurity: "max-age=31536000; includeSubDomains",
				ContentSecurityPolicy:   "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; img-src 'self' data:",
				ReferrerPolicy:          "strict-origin-when-cross-origin",
			},
		},

		Logging: LoggingConfig{
			Level:      getLogLevel("LOG_LEVEL", slog.LevelInfo),
			Format:     getEnv("LOG_FORMAT", getDefaultLogFormat()),
			AddSource:  getBoolEnv("LOG_ADD_SOURCE", false),
			TimeFormat: getEnv("LOG_TIME_FORMAT", time.RFC3339),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getIntEnv("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
			PoolSize: getIntEnv("REDIS_POOL_SIZE", 10),
		},

		GRPC: GRPCConfig{
			ServerPort:       getIntEnv("GRPC_SERVER_PORT", 9090),
			WorkerAddress:    getEnv("WORKER_GRPC_ADDRESS", "worker-server:9091"),
			TLSCertFile:      getEnv("TLS_CERT_FILE", "./certs/server.crt"),
			TLSKeyFile:       getEnv("TLS_KEY_FILE", "./certs/server.key"),
			TLSEnabled:       getBoolEnv("TLS_ENABLED", false), // Disabled by default for development
			APIServerAddress: getEnv("GRPC_API_SERVER_ADDRESS", "localhost:50051"),
			ServerName:       getEnv("GRPC_SERVER_NAME", ""),
		},

		Worker: WorkerConfig{
			HealthPort:         getIntEnv("WORKER_HEALTH_PORT", 8091),
			WorkerID:           getEnv("WORKER_ID", "worker-1"),
			WorkerName:         getEnv("WORKER_NAME", "EngLog Worker"),
			Version:            getEnv("WORKER_VERSION", "1.0.0"),
			OllamaURL:          getEnv("OLLAMA_URL", "http://localhost:11434"),
			MaxConcurrentTasks: getIntEnv("MAX_CONCURRENT_TASKS", 5),
		},
	}

	return cfg
}

// getDefaultLogFormat returns the default log format based on environment
func getDefaultLogFormat() string {
	if getEnv("APP_ENV", "development") == "production" {
		return "json"
	}
	return "text"
}

// getLogLevel parses log level from string
func getLogLevel(key string, defaultLevel slog.Level) slog.Level {
	levelStr := getEnv(key, "")
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return defaultLevel
	}
}

// getSliceEnv parses a comma-separated string into a slice
func getSliceEnv(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// Helper functions for environment variable parsing

// getEnv gets string environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets integer environment variable with default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv gets boolean environment variable with default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getDurationEnv gets duration environment variable with default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
