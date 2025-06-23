package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type PatientHandler struct {
	svc port.PatientService
}

func NewPatientHandler(svc port.PatientService) *PatientHandler {
	return &PatientHandler{
		svc: svc,
	}
}

// CreatePatient godoc
//
//	@Summary		Create a patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.CreatePatientReq	true	"Patient details"
