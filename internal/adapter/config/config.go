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
