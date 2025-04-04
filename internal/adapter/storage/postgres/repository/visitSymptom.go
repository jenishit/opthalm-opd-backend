package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

type VisitSymptomRepository struct {
	DB *postgres.DB
}

func NewVisitSymptomRepository(db *postgres.DB) *VisitSymptomRepository {
	return &VisitSymptomRepository{
		DB: db,
	}
}

func (vsr *VisitSymptomRepository) UpsertVisitSymptom(ctx context.Context, vr *domain.VisitSymptom) (*domain.CreateVSResponse, error) {
	query := `INSERT INTO
	VISIT_SYMPTOMS (
		VISIT_ID,
		FOREIGN_BODY_SENSATION,
		REDNESS,
		SWELLING,
		DIMINISHED_DISTANCE_VISION,
		DIMINISHED_NEAR_VISION,
		BLURRING,
		HEADACHE,
		EYE_PAIN,
		WATERY_EYES,
		DISCHARGE,
		ITCHING,
		PHOTOPHOBIA,
		FLOATERS,
		FLASHES,
		DOUBLE_VISION,
		RECENT_FEVER,
		COMPLAINT_DURATION,
		CREATED_BY,
		UPDATED_BY
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20
	)
ON CONFLICT (VISIT_ID) DO UPDATE
SET
