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
