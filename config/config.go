package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppName     string
	AppPort     string
	AppTimezone string
	Location    *time.Location

	DatabaseURL   string
	DBMaxOpenConn int
	DBMaxIdleConn int

	JWTSecret   string
	JWTIssuer   string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
	CORSOrigins []string

	R2 R2Config

	Seed SeedConfig
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string
	PublicBaseURL   string
}

type SeedConfig struct {
	DeveloperName     string
	DeveloperUsername string
	DeveloperEmail    string
	DeveloperPassword string
}

func (r R2Config) Enabled() bool {
	return r.AccountID != "" && r.AccessKeyID != "" && r.SecretAccessKey != "" && r.Bucket != ""
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	loc, err := time.LoadLocation(env("APP_TIMEZONE", "Asia/Jakarta"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_TIMEZONE: %w", err)
	}

	accessTTL, err := time.ParseDuration(env("JWT_ACCESS_TTL", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}

	refreshTTL, err := time.ParseDuration(env("JWT_REFRESH_TTL", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}

	dbMaxOpen, err := envInt("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return nil, err
	}
	dbMaxIdle, err := envInt("DB_MAX_IDLE_CONNS", 10)
	if err != nil {
		return nil, err
	}
	if dbMaxIdle > dbMaxOpen {
		return nil, fmt.Errorf("DB_MAX_IDLE_CONNS must be less than or equal to DB_MAX_OPEN_CONNS")
	}

	appEnv := env("APP_ENV", "development")
	switch appEnv {
	case "development", "test", "staging", "production":
	default:
		return nil, fmt.Errorf("APP_ENV must be one of development, test, staging, production")
	}

	secret := env("JWT_SECRET", "")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if appEnv == "production" {
		for _, placeholder := range []string{"change-me", "replace-with", "dev-only", "ci-only"} {
			if strings.Contains(strings.ToLower(secret), placeholder) {
				return nil, fmt.Errorf("JWT_SECRET must not use a placeholder value in production")
			}
		}
	}

	r2 := R2Config{
		AccountID:       env("R2_ACCOUNT_ID", ""),
		AccessKeyID:     env("R2_ACCESS_KEY_ID", ""),
		SecretAccessKey: env("R2_SECRET_ACCESS_KEY", ""),
		Bucket:          env("R2_BUCKET", ""),
		Endpoint:        env("R2_ENDPOINT", ""),
		PublicBaseURL:   strings.TrimRight(env("R2_PUBLIC_BASE_URL", ""), "/"),
	}
	if r2.Endpoint == "" && r2.AccountID != "" {
		r2.Endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2.AccountID)
	}

	return &Config{
		AppEnv:        appEnv,
		AppName:       env("APP_NAME", "jiandco-api"),
		AppPort:       env("APP_PORT", "8080"),
		AppTimezone:   env("APP_TIMEZONE", "Asia/Jakarta"),
		Location:      loc,
		DatabaseURL:   env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/jiandco_db?sslmode=disable"),
		DBMaxOpenConn: dbMaxOpen,
		DBMaxIdleConn: dbMaxIdle,
		JWTSecret:     secret,
		JWTIssuer:     env("JWT_ISSUER", "jiandco-api"),
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
		CORSOrigins:   splitCSV(env("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		R2:            r2,
		Seed: SeedConfig{
			DeveloperName:     env("SEED_DEVELOPER_NAME", "Developer"),
			DeveloperUsername: env("SEED_DEVELOPER_USERNAME", "developer"),
			DeveloperEmail:    env("SEED_DEVELOPER_EMAIL", "developer@jiandco.local"),
			DeveloperPassword: env("SEED_DEVELOPER_PASSWORD", ""),
		},
	}, nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) Addr() string {
	return ":" + c.AppPort
}

func (c *Config) AccessTTLSeconds() int {
	return int(c.AccessTTL.Seconds())
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) (int, error) {
	raw := env(key, strconv.Itoa(fallback))
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return value, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
