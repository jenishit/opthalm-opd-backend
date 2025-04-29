package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Container contains environment variables for the application, database, cache, token, and http server
type (
	Container struct {
		App     *App
		Token   *Token
		Refresh *Refresh

		Session *Session
		DB      *DB
		HTTP    *HTTP

		Redis *RedisConfig
		SMTP  *SMTP
	}
	// App contains all the environment variables for the application
	App struct {
		Name string
		Env  string
	}
	// Token contains all the environment variables for the token service
	Token struct {
		Secret   string
		Duration string
		Refresh  *Refresh
	}
	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
		DB       int
		Prefix   string
	}

	RedisConfig struct {
		Addr     string
		Password string
		DB       int
	}
	Session struct {
		Driver string
		TTL    time.Duration
		Redis  *Redis
	}
	Cache struct {
		Enabled    bool
		Driver     string
		DefaultTTL time.Duration
		Redis      *Redis
	}
	Refresh struct {
		Duration string
	}
	// Database contains all the environment variables for the database
	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env                string
		URL                string
		Port               string
		AllowedOrigins     string
		UseFunctionURLCORS bool
	}
	// SMTP contains the environment variables for outbound email. When Host
	// is empty, main.go wires a log-only EmailSender instead of net/smtp, so
	// password-reset/verification flows still work end-to-end without a
	// mail server configured — just without real delivery.
	SMTP struct {
		Host     string
		Port     string
		Username string
		Password string
		From     string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	token := &Token{
		Duration: os.Getenv("TOKEN_DURATION"),
		Secret:   os.Getenv("TOKEN_SECRET"),
	}

	refreshToken := &Refresh{
		Duration: os.Getenv("REFRESH_TOKEN_DURATION"),
	}

	redisAddr := envOrDefault(os.Getenv("REDIS_ADDR"), "")

	redisPassword := envOrDefault(os.Getenv("REDIS_PASSWORD"), "")

	session := &Session{
		Driver: strings.ToLower(envOrDefault(os.Getenv("SESSION_DRIVER"), "")),
		TTL:    parseDurationOrDefault(os.Getenv("SESSION_TTL"), 24*time.Hour),
		Redis: &Redis{
			Addr:     envOrDefault(os.Getenv("SESSION_REDIS_ADDR"), redisAddr),
			Password: envOrDefault(os.Getenv("SESSION_REDIS_PASSWORD"), redisPassword),
			DB:       parseIntOrDefault(os.Getenv("SESSION_REDIS_DB"), 0),
			Prefix:   envOrDefault(os.Getenv("SESSION_REDIS_PREFIX"), "session"),
		},
	}
	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	redisCfg := &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       parseIntOrDefault(os.Getenv("REDIS_DB"), 0),
	}

	http := &HTTP{
		Env:                os.Getenv("APP_ENV"),
		URL:                os.Getenv("HTTP_URL"),
		Port:               os.Getenv("HTTP_PORT"),
		AllowedOrigins:     os.Getenv("HTTP_ALLOWED_ORIGINS"),
		UseFunctionURLCORS: parseBool(os.Getenv("HTTP_USE_FUNCTION_URL_CORS"), false),
	}

	smtp := &SMTP{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     envOrDefault(os.Getenv("SMTP_PORT"), "587"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     envOrDefault(os.Getenv("SMTP_FROM"), "no-reply@example.com"),
	}

	return &Container{
		app,
		token,
		refreshToken,
		session,
		db,
		http,
		redisCfg,
		smtp,
	}, nil
}

func parseBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}
	switch strings.ToLower((value)) {
	case "1", "true", "yes", "on":
		return true

	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func envOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func parseDurationOrDefault(value string, fallback time.Duration) time.Duration {
	if value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return d
}

func parseIntOrDefault(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}

func parseFloatOrDefault(value string, fallback float64) float64 {
	if value == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return f
}

func resolveDatabaseConnection(secret *DB) string {
	return resolveDatabaseField(
		secretValue(secret, func(s *DB) string { return s.Connection }),
		os.Getenv("DB_CONNECTION"),
		"postgres",
	)
}

func resolveDatabaseHost(secret *DB) string {
	return resolveDatabaseField(
		secretValue(secret, func(s *DB) string { return s.Host }),
		envOrDefault(os.Getenv("DB_HOST"), os.Getenv("DATABASE_HOST")),
		"",
	)
}

func resolveDatabasePort(secret *DB) string {
	return resolveDatabaseField(
		secretValue(secret, func(s *DB) string { return s.Port }),
		os.Getenv("DB_PORT"),
		"5432",
	)
}

func resolveDatabaseUser(secret *DB) string {
	return resolveDatabaseField(
		secretValue(secret, func(s *DB) string { return s.User }),
		envOrDefault(os.Getenv("DB_USER"), os.Getenv("DATABASE_USER")),
		"",
	)
}

func resolveDatabaseName(secret *DB) string {
	return resolveDatabaseField(
		secretValue(secret, func(s *DB) string { return s.Name }),
		envOrDefault(os.Getenv("DB_NAME"), os.Getenv("DATABASE_NAME")),
		"",
	)
}

func resolveDatabaseSSLMode() string {
	if IsLambdaRuntime() {
		return envOrDefault(os.Getenv("SSL_MODE"), "require")
	}
	return os.Getenv("SSL_MODE")
}

func resolveDatabaseField(secretValue string, envValue string, fallback string) string {
	if IsLambdaRuntime() {
		return envOrDefault(secretValue, envOrDefault(envValue, fallback))
	}
	return envOrDefault(envValue, envOrDefault(secretValue, fallback))
}

