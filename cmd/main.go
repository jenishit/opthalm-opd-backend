// Package main is the Opthalmic Management System API server.
//
//	@title			Opthalmic Management System API
//	@version		1.0
//	@description	Backend API for a multi-tenant (SaaS) ophthalmology clinic
//	@description	management system: patients, visits, billing/POS, inventory,
//	@description	lab jobs, reports, and optical calculators.
//	@description
//	@description	Every route under /api (other than /api/auth/*, /api/signup,
//	@description	and the /api/auth/password-reset|email endpoints) requires a
//	@description	Bearer access token, and most additionally require the
//	@description	caller's clinic to have an active subscription.
//
//	@contact.name	API Support
//
//	@license.name	Proprietary
//
//	@host		localhost:8082
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and the access token, e.g. "Bearer eyJhbGci...".
package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"os"
