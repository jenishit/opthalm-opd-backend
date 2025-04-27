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

