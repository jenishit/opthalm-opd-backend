package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
)

// ─── Visit Symptoms ─────────────────────────────────────────────────────────

type UpsertVisitSymptomReq struct {
	ForeignBodySensation     string `json:"foreign_body_sensation" binding:"omitempty,oneof=none OD OS OU"`
	Redness                  string `json:"redness" binding:"omitempty,oneof=none OD OS OU"`
	Swelling                 string `json:"swelling" binding:"omitempty,oneof=none OD OS OU"`
	DiminishedDistanceVision string `json:"diminished_distance_vision" binding:"omitempty,oneof=none OD OS OU"`
	DiminishedNearVision     string `json:"diminished_near_vision" binding:"omitempty,oneof=none OD OS OU"`
	Blurring                 string `json:"blurring" binding:"omitempty,oneof=none OD OS OU"`
	Headache                 string `json:"headache" binding:"omitempty,oneof=none OD OS OU"`
	EyePain                  string `json:"eye_pain" binding:"omitempty,oneof=none OD OS OU"`
	WateryEyes               string `json:"watery_eyes" binding:"omitempty,oneof=none OD OS OU"`
	Discharge                string `json:"discharge" binding:"omitempty,oneof=none OD OS OU"`
	Itching                  string `json:"itching" binding:"omitempty,oneof=none OD OS OU"`
	Photophobia              string `json:"photophobia" binding:"omitempty,oneof=none OD OS OU"`
	Floaters                 string `json:"floaters" binding:"omitempty,oneof=none OD OS OU"`
	Flashes                  string `json:"flashes" binding:"omitempty,oneof=none OD OS OU"`
	DoubleVision             string `json:"double_vision" binding:"omitempty,oneof=none OD OS OU"`
	RecentFever              bool   `json:"recent_fever"`
}

type VisitSymptomResponse struct {
	VisitID                  uuid.UUID `json:"visit_id"`
	ForeignBodySensation     string    `json:"foreign_body_sensation"`
	Redness                  string    `json:"redness"`
	Swelling                 string    `json:"swelling"`
	DiminishedDistanceVision string    `json:"diminished_distance_vision"`
	DiminishedNearVision     string    `json:"diminished_near_vision"`
	Blurring                 string    `json:"blurring"`
	Headache                 string    `json:"headache"`
	EyePain                  string    `json:"eye_pain"`
	WateryEyes               string    `json:"watery_eyes"`
	Discharge                string    `json:"discharge"`
	Itching                  string    `json:"itching"`
	Photophobia              string    `json:"photophobia"`
	Floaters                 string    `json:"floaters"`
	Flashes                  string    `json:"flashes"`
	DoubleVision             string    `json:"double_vision"`
	RecentFever              bool      `json:"recent_fever"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func VisitSymptomRes(vs *domain.VisitSymptom) *VisitSymptomResponse {
	if vs == nil {
		return nil
	}

	return &VisitSymptomResponse{
		VisitID:                  vs.VisitID,
		ForeignBodySensation:     strOrDefault(vs.ForeignBodySensation),
		Redness:                  strOrDefault(vs.Redness),
		Swelling:                 strOrDefault(vs.Swelling),
		DiminishedDistanceVision: strOrDefault(vs.DiminishedDistanceVision),
		DiminishedNearVision:     strOrDefault(vs.DiminishedNearVision),
		Blurring:                 strOrDefault(vs.Blurring),
		Headache:                 strOrDefault(vs.Headache),
		EyePain:                  strOrDefault(vs.EyePain),
		WateryEyes:               strOrDefault(vs.WateryEyes),
		Discharge:                strOrDefault(vs.Discharge),
		Itching:                  strOrDefault(vs.Itching),
		Photophobia:              strOrDefault(vs.PhotoPhobia),
		Floaters:                 strOrDefault(vs.Floaters),
		Flashes:                  strOrDefault(vs.Flashes),
		DoubleVision:             strOrDefault(vs.DoubleVision),
		RecentFever:              boolOrDefault(vs.RecentFever),
		UpdatedAt:                vs.UpdatedAt,
	}
}

// ─── Examination Findings ──────────────────────────────────────────────────

type UpsertExaminationFindingsReq struct {
	IOPOD             *float64 `json:"iop_od"`
	IOPOS             *float64 `json:"iop_os"`
	IOPTime           *string  `json:"iop_time"`
	BP                *string  `json:"bp"`
	Convergence       *string  `json:"convergence"`
	CoverTest         *string  `json:"cover_test"`
	EOM               *string  `json:"eom"`
	AnteriorSegmentOD *string  `json:"anterior_segment_od"`
	AnteriorSegmentOS *string  `json:"anterior_segment_os"`
	FundusOD          *string  `json:"fundus_od"`
	FundusOS          *string  `json:"fundus_os"`
}

type ExaminationFindingsResponse struct {
	ID                uuid.UUID `json:"id"`
	VisitID           uuid.UUID `json:"visit_id"`
	IOPOD             *float64  `json:"iop_od,omitempty"`
	IOPOS             *float64  `json:"iop_os,omitempty"`
	IOPTime           *string   `json:"iop_time,omitempty"`
	BP                *string   `json:"bp,omitempty"`
	Convergence       *string   `json:"convergence,omitempty"`
	CoverTest         *string   `json:"cover_test,omitempty"`
	EOM               *string   `json:"eom,omitempty"`
	AnteriorSegmentOD *string   `json:"anterior_segment_od,omitempty"`
	AnteriorSegmentOS *string   `json:"anterior_segment_os,omitempty"`
	FundusOD          *string   `json:"fundus_od,omitempty"`
	FundusOS          *string   `json:"fundus_os,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ExaminationFindingsRes(ef *domain.ExaminationFindings) *ExaminationFindingsResponse {
	if ef == nil {
		return nil
	}

	return &ExaminationFindingsResponse{
		ID:                ef.ID,
		VisitID:           ef.VisitID,
		IOPOD:             ef.IOPOD,
		IOPOS:             ef.IOPOS,
		IOPTime:           ef.IOPTime,
		BP:                ef.BP,
		Convergence:       ef.Convergence,
		CoverTest:         ef.CoverTest,
		EOM:               ef.EOM,
		AnteriorSegmentOD: ef.AnteriorSegmentOD,
		AnteriorSegmentOS: ef.AnteriorSegmentOS,
		FundusOD:          ef.FundusOD,
		FundusOS:          ef.FundusOS,
		UpdatedAt:         ef.UpdatedAt,
	}
}

// ─── Refraction Readings ────────────────────────────────────────────────────

type AddRefractionReadingReq struct {
	Eye      string   `json:"eye" binding:"required,oneof=OD OS"`
	Stage    string   `json:"stage" binding:"required,oneof=unaided old_rx pinhole retinoscopy subjective_acceptance final_rx"`
	Sphere   *float64 `json:"sphere"`
	Cylinder *float64 `json:"cylinder"`
	Axis     *int     `json:"axis" binding:"omitempty,min=0,max=180"`
	AddPower *float64 `json:"add_power"`
	VA       *string  `json:"va"`
	PD       *float64 `json:"pd"`
	LensType *string  `json:"lens_type" binding:"omitempty,oneof=single_vision bifocal progressive contact_lens none"`
}

type RefractionReadingResponse struct {
	ID        uuid.UUID `json:"id"`
	VisitID   uuid.UUID `json:"visit_id"`
	Eye       string    `json:"eye"`
	Stage     string    `json:"stage"`
	Sphere    *float64  `json:"sphere,omitempty"`
	Cylinder  *float64  `json:"cylinder,omitempty"`
	Axis      *int      `json:"axis,omitempty"`
	AddPower  *float64  `json:"add_power,omitempty"`
	VA        *string   `json:"va,omitempty"`
	PD        *float64  `json:"pd,omitempty"`
	LensType  *string   `json:"lens_type,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func RefractionReadingRes(r *domain.RefractionReading) *RefractionReadingResponse {
	if r == nil {
		return nil
	}

	return &RefractionReadingResponse{
		ID:        r.ID,
		VisitID:   r.VisitID,
		Eye:       r.Eye,
		Stage:     r.Stage,
		Sphere:    r.Sphere,
		Cylinder:  r.Cylinder,
		Axis:      r.Axis,
		AddPower:  r.AddPower,
		VA:        r.VA,
		PD:        r.PD,
		LensType:  r.LensType,
		CreatedAt: r.CreatedAt,
	}
}

func RefractionReadingResList(readings []*domain.RefractionReading) []*RefractionReadingResponse {
	rsp := make([]*RefractionReadingResponse, 0, len(readings))
	for _, r := range readings {
		rsp = append(rsp, RefractionReadingRes(r))
	}
	return rsp
}

// ─── Visit Diagnoses ─────────────────────────────────────────────────────────

type VisitDiagnosisReq struct {
	DiagnosisID   uuid.UUID `json:"diagnosis_id" binding:"required"`
	DiagnosisName *string   `json:"diagnosis_name"`
	Icd10Code     *string   `json:"icd10_code"`
	Laterality    string    `json:"laterality" binding:"required,oneof=OD OS OU none"`
	Status        *string   `json:"status" binding:"omitempty,oneof=provisional confirmed resolved"`
}

type ReplaceVisitDiagnosesReq struct {
	Diagnoses []VisitDiagnosisReq `json:"diagnoses" binding:"dive"`
}

type VisitDiagnosisResponse struct {
	ID            uuid.UUID `json:"id"`
	VisitID       uuid.UUID `json:"visit_id"`
	DiagnosisID   uuid.UUID `json:"diagnosis_id"`
	DiagnosisName *string   `json:"diagnosis_name,omitempty"`
	Icd10Code     *string   `json:"icd10_code,omitempty"`
	Laterality    string    `json:"laterality"`
	Status        *string   `json:"status,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func VisitDiagnosisRes(d *domain.VisitDiagnosis) *VisitDiagnosisResponse {
	if d == nil {
		return nil
	}

	return &VisitDiagnosisResponse{
		ID:            d.ID,
		VisitID:       d.VisitID,
		DiagnosisID:   d.DiagnosisID,
		DiagnosisName: d.DiagnosisName,
		Icd10Code:     d.Icd10Code,
		Laterality:    d.Laterality,
		Status:        d.Status,
		CreatedAt:     d.CreatedAt,
	}
}

func VisitDiagnosisResList(diagnoses []*domain.VisitDiagnosis) []*VisitDiagnosisResponse {
	rsp := make([]*VisitDiagnosisResponse, 0, len(diagnoses))
	for _, d := range diagnoses {
		rsp = append(rsp, VisitDiagnosisRes(d))
	}
	return rsp
}

// ─── Treatments ──────────────────────────────────────────────────────────────

type TreatmentReq struct {
	Type          string     `json:"type" binding:"required,oneof=medicine procedure advice"`
	MedicineID    *uuid.UUID `json:"medicine_id"`
	MedicineName  *string    `json:"medicine_name"`
	ProcedureName *string    `json:"procedure_name"`
	AdviceText    *string    `json:"advice_text"`
	Route         *string    `json:"route"`
	Dosage        *string    `json:"dosage"`
	Frequency     *string    `json:"frequency"`
	Duration      *string    `json:"duration"`
}

type ReplaceTreatmentsReq struct {
	Treatments []TreatmentReq `json:"treatments" binding:"dive"`
}

type TreatmentResponse struct {
	ID            uuid.UUID  `json:"id"`
	VisitID       uuid.UUID  `json:"visit_id"`
	Type          string     `json:"type"`
	MedicineID    *uuid.UUID `json:"medicine_id,omitempty"`
	MedicineName  *string    `json:"medicine_name,omitempty"`
	ProcedureName *string    `json:"procedure_name,omitempty"`
	AdviceText    *string    `json:"advice_text,omitempty"`
	Route         *string    `json:"route,omitempty"`
	Dosage        *string    `json:"dosage,omitempty"`
	Frequency     *string    `json:"frequency,omitempty"`
	Duration      *string    `json:"duration,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func TreatmentRes(t *domain.Treatment) *TreatmentResponse {
	if t == nil {
		return nil
	}

	return &TreatmentResponse{
		ID:            t.ID,
		VisitID:       t.VisitID,
		Type:          t.Type,
		MedicineID:    t.MedicineID,
		MedicineName:  t.MedicineName,
		ProcedureName: t.ProcedureName,
		AdviceText:    t.AdviceText,
		Route:         t.Route,
		Dosage:        t.Dosage,
		Frequency:     t.Frequency,
		Duration:      t.Duration,
		CreatedAt:     t.CreatedAt,
	}
}

func TreatmentResList(treatments []*domain.Treatment) []*TreatmentResponse {
	rsp := make([]*TreatmentResponse, 0, len(treatments))
	for _, t := range treatments {
		rsp = append(rsp, TreatmentRes(t))
	}
	return rsp
}

// ─── Follow Up ───────────────────────────────────────────────────────────────

type UpsertFollowUpReq struct {
	FollowUpDate string  `json:"follow_up_date" binding:"required,datetime=2006-01-02"`
	Reason       *string `json:"reason"`
}

type FollowUpResponse struct {
	ID             uuid.UUID  `json:"id"`
	VisitID        uuid.UUID  `json:"visit_id"`
	FollowUpDate   string     `json:"follow_up_date"`
	Reason         *string    `json:"reason,omitempty"`
	ReminderSentAt *time.Time `json:"reminder_sent_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func FollowUpRes(f *domain.FollowUp) *FollowUpResponse {
	if f == nil {
		return nil
	}

	return &FollowUpResponse{
		ID:             f.ID,
		VisitID:        f.VisitID,
		FollowUpDate:   f.FollowUpDate,
		Reason:         f.Reason,
		ReminderSentAt: f.ReminderSentAt,
		UpdatedAt:      f.UpdatedAt,
	}
}

// ─── Full clinical chart (GET /visit/:id/clinical) ─────────────────────────

type VisitClinicalResponse struct {
	Symptoms    *VisitSymptomResponse        `json:"symptoms,omitempty"`
	Examination *ExaminationFindingsResponse `json:"examination,omitempty"`
	Readings    []*RefractionReadingResponse `json:"readings"`
	Diagnoses   []*VisitDiagnosisResponse    `json:"diagnoses"`
	Treatments  []*TreatmentResponse         `json:"treatments"`
	FollowUp    *FollowUpResponse            `json:"follow_up,omitempty"`
}

func strOrDefault(s *string) string {
	if s == nil {
		return "none"
	}
	return *s
}

func boolOrDefault(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
