// Package testutil provides an integration-test harness that wires up the
// full application (exactly as cmd/main.go does) against a real Postgres
// database and exposes it via httptest, plus small helpers for the
// login/create-patient/etc. boilerplate every test needs.
//
// It deliberately does NOT use mocked repositories: the bugs this suite is
// meant to catch (pgx type-scanning mismatches, missing NOT NULL columns,
// broken SQL fragments) only show up against a real database.
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	auth "github.com/jenish-brainztechs/go-backend/internal/adapter/auth/jwt"
	redisadapter "github.com/jenish-brainztechs/go-backend/internal/adapter/cache/redis"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	httpadapter "github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres/repository"
	"github.com/jenish-brainztechs/go-backend/internal/core/services"
