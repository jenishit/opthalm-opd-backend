package services

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type PatientService struct {
	repo port.PatientRepository
}

func NewPatientService(pr port.PatientRepository) *PatientService {
	return &PatientService{
		repo: pr,
	}
}

func (ps *PatientService) CreatePatient(ctx *gin.Context)
func (ps *PatientService) GetPatientByID(ctx context.Context)
func (ps *PatientService) GetPatients(ctx *gin.Context)
func (ps *PatientService) UpdatePatientByID(ctx *gin.Context)
func (ps *PatientService) DeletePatientByID(ctx *gin.Context)
