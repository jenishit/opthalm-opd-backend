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
