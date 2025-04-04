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
	FOREIGN_BODY_SENSATION = EXCLUDED.FOREIGN_BODY_SENSATION,
	REDNESS = EXCLUDED.REDNESS,
	SWELLING = EXCLUDED.SWELLING,
	DIMINISHED_DISTANCE_VISION = EXCLUDED.DIMINISHED_DISTANCE_VISION,
	DIMINISHED_NEAR_VISION = EXCLUDED.DIMINISHED_NEAR_VISION,
	BLURRING = EXCLUDED.BLURRING,
	HEADACHE = EXCLUDED.HEADACHE,
	EYE_PAIN = EXCLUDED.EYE_PAIN,
	WATERY_EYES = EXCLUDED.WATERY_EYES,
	DISCHARGE = EXCLUDED.DISCHARGE,
	ITCHING = EXCLUDED.ITCHING,
	PHOTOPHOBIA = EXCLUDED.PHOTOPHOBIA,
	FLOATERS = EXCLUDED.FLOATERS,
	FLASHES = EXCLUDED.FLASHES,
	DOUBLE_VISION = EXCLUDED.DOUBLE_VISION,
	RECENT_FEVER = EXCLUDED.RECENT_FEVER,
	COMPLAINT_DURATION = EXCLUDED.COMPLAINT_DURATION,
	UPDATED_BY = EXCLUDED.UPDATED_BY,
	UPDATED_AT = NOW()
RETURNING
	VISIT_ID,
	CREATED_BY,
	UPDATED_BY,
	UPDATED_AT
	`

	var vs domain.CreateVSResponse

	err := vsr.DB.QueryRow(ctx, query).Scan(
		&vs.VisitID,
		&vs.CreatedBy,
		&vs.UpdatedBy,
		&vs.UpdatedAt,
	)

	if err != nil {
