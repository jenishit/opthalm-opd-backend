package domain

import (
	"time"

	"github.com/google/uuid"
)

type SalesSummary struct {
	Period       string
	InvoiceCount int
	Subtotal     float64
	Discount     float64
	Tax          float64
	TotalSales   float64
	TotalPaid    float64
	TotalDue     float64
}

type PatientDue struct {
	PatientID    uuid.UUID
	PatientName  string
	PatientPhone string
	InvoiceNo    string
	InvoiceID    uuid.UUID
	TotalAmount  float64
	PaidAmount   float64
	DueAmount    float64
	CreatedAt    time.Time
}

type VendorDue struct {
	VendorID     uuid.UUID
	VendorName   string
	PurchaseID   uuid.UUID
	TotalAmount  float64
	PaidAmount   float64
	DueAmount    float64
	PurchaseDate time.Time
}

type InventoryValuation struct {
	TotalItems     int
	TotalUnits     int
	TotalCostValue float64
	TotalSellValue float64
}

type VisitsSummary struct {
	Period     string
	VisitCount int
}

type DoctorVisitsSummary struct {
	DoctorID   uuid.UUID
	DoctorName string
	VisitCount int
}

// CategoryRevenue is one row of RevenueByCategory — Category is a display
// label ("Frames", "Lenses", ...), not the raw inventory_items.category DB
// value.
type CategoryRevenue struct {
	Category string
	Amount   float64
}

type RevenueByCategory struct {
	Period     string
	Categories []CategoryRevenue
}

// ClinicalActivitySummary counts, over a date range: new patients
// registered; visits with examination findings recorded
// (ComprehensiveExams); visits with a refraction reading recorded
// (RefractionVisits); and follow-ups due in the period (FollowUpVisits —
// scheduled follow-up date, not necessarily a visit that has happened yet,
// since nothing in the schema links a later visit back to the follow-up
// that scheduled it).
type ClinicalActivitySummary struct {
	Period             string
	NewPatients        int
	ComprehensiveExams int
	RefractionVisits   int
	FollowUpVisits     int
}
