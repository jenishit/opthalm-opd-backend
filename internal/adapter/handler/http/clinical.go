package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/handler/http/dto"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jenishit/opthalm-opd-backend/internal/core/port"
)

// ClinicalHandler serves the per-visit clinical chart: symptoms, examination
// findings, refraction readings, diagnoses, treatments, and follow-up —
// everything a doctor charts during an OPD visit besides the visit record
// itself (that's VisitHandler).
type ClinicalHandler struct {
	symptomSvc    port.VisitSymptomService
	examSvc       port.ExaminationFindingsService
	refractionSvc port.RefractionReadingService
	diagnosisSvc  port.VisitDiagnosisService
	treatmentSvc  port.TreatmentService
	followUpSvc   port.FollowUpService
}

func NewClinicalHandler(
	symptomSvc port.VisitSymptomService,
	examSvc port.ExaminationFindingsService,
	refractionSvc port.RefractionReadingService,
	diagnosisSvc port.VisitDiagnosisService,
	treatmentSvc port.TreatmentService,
	followUpSvc port.FollowUpService,
) *ClinicalHandler {
	return &ClinicalHandler{
		symptomSvc:    symptomSvc,
		examSvc:       examSvc,
		refractionSvc: refractionSvc,
		diagnosisSvc:  diagnosisSvc,
		treatmentSvc:  treatmentSvc,
		followUpSvc:   followUpSvc,
	}
}

// GetVisitClinical godoc
//
//	@Summary		Get a visit's full clinical chart
//	@Description	Symptoms, examination findings, refraction readings, diagnoses, treatments, and follow-up for one visit, in a single call.
//	@Tags			clinical
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Visit ID"
//	@Success		200	{object}	response{data=dto.VisitClinicalResponse}
//	@Router			/visit/{id}/clinical [get]
func (h *ClinicalHandler) GetVisitClinical(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	rsp := &dto.VisitClinicalResponse{}

	symptoms, err := h.symptomSvc.GetVisitSymptomByVisitID(ctx, user.ClinicID, visitID)
	if err != nil && err != domain.ErrDataNotFound {
		handleError(ctx, err)
		return
	}
	rsp.Symptoms = dto.VisitSymptomRes(symptoms)

	exam, err := h.examSvc.GetExaminationFindingsByVisitID(ctx, user.ClinicID, visitID)
	if err != nil && err != domain.ErrDataNotFound {
		handleError(ctx, err)
		return
	}
	rsp.Examination = dto.ExaminationFindingsRes(exam)

	readings, err := h.refractionSvc.ListReadingsByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp.Readings = dto.RefractionReadingResList(readings)

	diagnoses, err := h.diagnosisSvc.ListVisitDiagnosesByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp.Diagnoses = dto.VisitDiagnosisResList(diagnoses)

	treatments, err := h.treatmentSvc.ListTreatmentsByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp.Treatments = dto.TreatmentResList(treatments)

	followUp, err := h.followUpSvc.GetFollowUpByVisitID(ctx, user.ClinicID, visitID)
	if err != nil && err != domain.ErrDataNotFound {
		handleError(ctx, err)
		return
	}
	rsp.FollowUp = dto.FollowUpRes(followUp)

	handleSuccess(ctx, rsp)
}

// UpsertVisitSymptoms godoc
//
//	@Summary		Record a visit's symptoms
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Visit ID"
//	@Param			request	body		dto.UpsertVisitSymptomReq	true	"Symptoms"
//	@Success		200		{object}	response{data=dto.VisitSymptomResponse}
//	@Router			/visit/{id}/symptoms [put]
func (h *ClinicalHandler) UpsertVisitSymptoms(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.UpsertVisitSymptomReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	symptom := &domain.VisitSymptom{
		VisitID:                  visitID,
		ForeignBodySensation:     emptyToNone(req.ForeignBodySensation),
		Redness:                  emptyToNone(req.Redness),
		Swelling:                 emptyToNone(req.Swelling),
		DiminishedDistanceVision: emptyToNone(req.DiminishedDistanceVision),
		DiminishedNearVision:     emptyToNone(req.DiminishedNearVision),
		Blurring:                 emptyToNone(req.Blurring),
		Headache:                 emptyToNone(req.Headache),
		EyePain:                  emptyToNone(req.EyePain),
		WateryEyes:               emptyToNone(req.WateryEyes),
		Discharge:                emptyToNone(req.Discharge),
		Itching:                  emptyToNone(req.Itching),
		PhotoPhobia:              emptyToNone(req.Photophobia),
		Floaters:                 emptyToNone(req.Floaters),
		Flashes:                  emptyToNone(req.Flashes),
		DoubleVision:             emptyToNone(req.DoubleVision),
		RecentFever:              &req.RecentFever,
		CreatedBy:                user.UserId,
		UpdatedBy:                user.UserId,
	}

	if _, err := h.symptomSvc.UpsertVisitSymptom(ctx, user.ClinicID, symptom); err != nil {
		handleError(ctx, err)
		return
	}

	saved, err := h.symptomSvc.GetVisitSymptomByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VisitSymptomRes(saved))
}

// UpsertExamination godoc
//
//	@Summary		Record a visit's examination findings
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string								true	"Visit ID"
//	@Param			request	body		dto.UpsertExaminationFindingsReq	true	"Examination findings"
//	@Success		200		{object}	response{data=dto.ExaminationFindingsResponse}
//	@Router			/visit/{id}/examination [put]
func (h *ClinicalHandler) UpsertExamination(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.UpsertExaminationFindingsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	exam := &domain.ExaminationFindings{
		VisitID:           visitID,
		IOPOD:             req.IOPOD,
		IOPOS:             req.IOPOS,
		IOPTime:           req.IOPTime,
		BP:                req.BP,
		Convergence:       req.Convergence,
		CoverTest:         req.CoverTest,
		EOM:               req.EOM,
		AnteriorSegmentOD: req.AnteriorSegmentOD,
		AnteriorSegmentOS: req.AnteriorSegmentOS,
		FundusOD:          req.FundusOD,
		FundusOS:          req.FundusOS,
		CreatedBy:         user.UserId,
		UpdatedBy:         user.UserId,
	}

	saved, err := h.examSvc.UpsertExaminationFindings(ctx, user.ClinicID, exam)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.ExaminationFindingsRes(saved))
}

// ListReadings godoc
//
//	@Summary		List a visit's refraction readings
//	@Tags			clinical
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Visit ID"
//	@Success		200	{object}	response{data=[]dto.RefractionReadingResponse}
//	@Router			/visit/{id}/readings [get]
func (h *ClinicalHandler) ListReadings(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	readings, err := h.refractionSvc.ListReadingsByVisitID(ctx, user.ClinicID, visitID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.RefractionReadingResList(readings))
}

// AddReading godoc
//
//	@Summary		Record a refraction reading for one eye and stage
//	@Description	Re-submitting the same eye and stage for a visit replaces the earlier reading.
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Visit ID"
//	@Param			request	body		dto.AddRefractionReadingReq	true	"Reading"
//	@Success		200		{object}	response{data=dto.RefractionReadingResponse}
//	@Router			/visit/{id}/readings [post]
func (h *ClinicalHandler) AddReading(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.AddRefractionReadingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	reading := &domain.RefractionReading{
		VisitID:   visitID,
		Eye:       req.Eye,
		Stage:     req.Stage,
		Sphere:    req.Sphere,
		Cylinder:  req.Cylinder,
		Axis:      req.Axis,
		AddPower:  req.AddPower,
		VA:        req.VA,
		PD:        req.PD,
		LensType:  req.LensType,
		CreatedBy: user.UserId,
	}

	saved, err := h.refractionSvc.AddReading(ctx, user.ClinicID, reading)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.RefractionReadingRes(saved))
}

// ReplaceDiagnoses godoc
//
//	@Summary		Replace a visit's diagnoses
//	@Description	Replaces the full set of diagnoses recorded for the visit.
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Visit ID"
//	@Param			request	body		dto.ReplaceVisitDiagnosesReq	true	"Diagnoses"
//	@Success		200		{object}	response{data=[]dto.VisitDiagnosisResponse}
//	@Router			/visit/{id}/diagnoses [put]
func (h *ClinicalHandler) ReplaceDiagnoses(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.ReplaceVisitDiagnosesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	diagnoses := make([]*domain.VisitDiagnosis, 0, len(req.Diagnoses))
	for _, d := range req.Diagnoses {
		diagnoses = append(diagnoses, &domain.VisitDiagnosis{
			VisitID:       visitID,
			DiagnosisID:   d.DiagnosisID,
			DiagnosisName: d.DiagnosisName,
			Icd10Code:     d.Icd10Code,
			Laterality:    d.Laterality,
			Status:        d.Status,
			CreatedBy:     user.UserId,
		})
	}

	saved, err := h.diagnosisSvc.ReplaceVisitDiagnoses(ctx, user.ClinicID, visitID, diagnoses)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.VisitDiagnosisResList(saved))
}

// ReplaceTreatments godoc
//
//	@Summary		Replace a visit's treatments
//	@Description	Replaces the full set of treatments (medicines, procedures, advice) recorded for the visit.
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Visit ID"
//	@Param			request	body		dto.ReplaceTreatmentsReq	true	"Treatments"
//	@Success		200		{object}	response{data=[]dto.TreatmentResponse}
//	@Router			/visit/{id}/treatments [put]
func (h *ClinicalHandler) ReplaceTreatments(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.ReplaceTreatmentsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	treatments := make([]*domain.Treatment, 0, len(req.Treatments))
	for _, t := range req.Treatments {
		treatments = append(treatments, &domain.Treatment{
			VisitID:       visitID,
			Type:          t.Type,
			MedicineID:    t.MedicineID,
			MedicineName:  t.MedicineName,
			ProcedureName: t.ProcedureName,
			AdviceText:    t.AdviceText,
			Route:         t.Route,
			Dosage:        t.Dosage,
			Frequency:     t.Frequency,
			Duration:      t.Duration,
			CreatedBy:     user.UserId,
		})
	}

	saved, err := h.treatmentSvc.ReplaceTreatments(ctx, user.ClinicID, visitID, treatments)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.TreatmentResList(saved))
}

// UpsertFollowUp godoc
//
//	@Summary		Schedule a visit's follow-up
//	@Tags			clinical
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Visit ID"
//	@Param			request	body		dto.UpsertFollowUpReq	true	"Follow-up"
//	@Success		200		{object}	response{data=dto.FollowUpResponse}
//	@Router			/visit/{id}/follow-up [put]
func (h *ClinicalHandler) UpsertFollowUp(ctx *gin.Context) {
	visitID, ok := parseVisitID(ctx)
	if !ok {
		return
	}

	var req dto.UpsertFollowUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user, ok := currentUser(ctx)
	if !ok {
		return
	}

	followUp := &domain.FollowUp{
		VisitID:      visitID,
		FollowUpDate: req.FollowUpDate,
		Reason:       req.Reason,
		CreatedBy:    user.UserId,
		UpdatedBy:    user.UserId,
	}

	saved, err := h.followUpSvc.UpsertFollowUp(ctx, user.ClinicID, followUp)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, dto.FollowUpRes(saved))
}

// parseVisitID extracts and validates the ":id" path param shared by every
// clinical-chart route, writing the error response itself on failure.
func parseVisitID(ctx *gin.Context) (uuid.UUID, bool) {
	visitID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		handleError(ctx, domain.ErrInvalidUUID)
		return uuid.Nil, false
	}
	return visitID, true
}

// emptyToNone defaults an omitted laterality field to "none", matching the
// symptom chart's column defaults.
func emptyToNone(s string) *string {
	if s == "" {
		s = "none"
	}
	return &s
}
