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

	http := &HTTP{
		Env:                os.Getenv("APP_ENV"),
		URL:                os.Getenv("HTTP_URL"),
		Port:               os.Getenv("HTTP_PORT"),
		AllowedOrigins:     os.Getenv("HTTP_ALLOWED_ORIGINS"),
		UseFunctionURLCORS: parseBool(os.Getenv("HTTP_USE_FUNCTION_URL_CORS"), false),
	}

	return &Container{
		app,
		token,
		refreshToken,
		session,
		db,
		http,
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

func firstSecretString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := payload[key]
		if !ok || value == nil {
			continue
		}
		switch v := value.(type) {
		case string:
			if trimmed := strings.TrimSpace(v); trimmed != "" {
				return trimmed
			}
		case float64:
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		}
	}
	return ""
}

func secretValue(secret *DB, selector func(*DB) string) string {
	if secret == nil {
		return ""
	}
	return selector(secret)
}

func normalizeDatabaseConnection(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "postgresql":
		return "postgres"
	default:
		return strings.TrimSpace(value)
	}
}

func IsLambdaRuntime() bool {
	return parseBool(os.Getenv("IS_LAMBDA_RUNTIME"), false)
}

func validateRuntimeDBConfig(db *DB) error {
	if db.Connection == "" {
		return fmt.Errorf("missing DB_CONNECTION")
	}
	if db.Host == "" {
		return fmt.Errorf("missing DB_HOST or DATABASE_HOST")
	}
	if db.Port == "" {
		return fmt.Errorf("missing DB_PORT")
	}
	if db.User == "" {
		return fmt.Errorf("missing DB_USER or DATABASE_USER")
	}
	if db.Name == "" {
		return fmt.Errorf("missing DB_NAME or DATABASE_NAME")
	}
	if db.Password == "" {
		return fmt.Errorf("missing DB_PASSWORD/DATABASE_PASSWORD and DATABASE_SECRET_ARN")
	}

	return nil
}
