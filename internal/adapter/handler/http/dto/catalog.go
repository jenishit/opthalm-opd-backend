package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

// ─── Medicine ─────────────────────────────────────────────────────────────────

type UpdateMedicineReq struct {
	MedicineName *string `json:"medicine_name"`
	BrandName    *string `json:"brand_name"`
	Strength     *string `json:"strength"`
	Form         *string `json:"form"`
}

type MedicineResponse struct {
	ID           uuid.UUID `json:"id"`
	MedicineName string    `json:"medicine_name"`
	BrandName    *string   `json:"brand_name"`
	Strength     *string   `json:"strength"`
	Form         *string   `json:"form"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func MedicineRes(m *domain.Medicine) *MedicineResponse {
	return &MedicineResponse{
		ID:           m.ID,
		MedicineName: m.MedicineName,
		BrandName:    m.BrandName,
		Strength:     m.Strength,
		Form:         m.Form,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func MedicineResList(ms []*domain.Medicine) []*MedicineResponse {
	items := make([]*MedicineResponse, 0, len(ms))
	for _, m := range ms {
		items = append(items, MedicineRes(m))
	}
	return items
}

// ─── Diagnosis Catalog ────────────────────────────────────────────────────────

type UpdateDiagnosisCatalogReq struct {
	Name      *string `json:"name"`
	Icd10Code *string `json:"icd10_code"`
}

type DiagnosisCatalogResponse struct {
	ID        uuid.UUID `json:"id"`
	Icd10Code *string   `json:"icd10_code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DiagnosisCatalogRes(d *domain.DiagnosisCatalog) *DiagnosisCatalogResponse {
	return &DiagnosisCatalogResponse{
		ID:        d.ID,
		Icd10Code: d.Icd10Code,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func DiagnosisCatalogResList(ds []*domain.DiagnosisCatalog) []*DiagnosisCatalogResponse {
	items := make([]*DiagnosisCatalogResponse, 0, len(ds))
	for _, d := range ds {
		items = append(items, DiagnosisCatalogRes(d))
	}
	return items
}

// ─── History Condition ─────────────────────────────────────────────────────────

type UpdateHistoryConditionReq struct {
	Name *string `json:"name"`
}

type HistoryConditionResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func HistoryConditionRes(h *domain.HistoryCondition) *HistoryConditionResponse {
	return &HistoryConditionResponse{
		ID:        h.ID,
		Name:      h.Name,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func HistoryConditionResList(hs []*domain.HistoryCondition) []*HistoryConditionResponse {
	items := make([]*HistoryConditionResponse, 0, len(hs))
	for _, h := range hs {
		items = append(items, HistoryConditionRes(h))
	}
	return items
}
