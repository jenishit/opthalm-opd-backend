package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type CreateInvoiceItemReq struct {
	BundleID        *uuid.UUID `json:"bundle_id"`
	ItemType        string     `json:"item_type" binding:"required,oneof=frame lens coating contact_lens service other"`
	Description     string     `json:"description" binding:"required"`
	InventoryItemID *uuid.UUID `json:"inventory_item_id"`
	Quantity        int        `json:"quantity" binding:"required,min=1"`
	UnitPrice       float64    `json:"unit_price" binding:"required,min=0"`
	DiscountAmount  float64    `json:"discount_amount"`
}

type CreateInvoiceReq struct {
	PatientID      uuid.UUID              `json:"patient_id" binding:"required"`
	VisitID        *uuid.UUID             `json:"visit_id"`
	DiscountAmount float64                `json:"discount_amount"`
	TaxAmount      float64                `json:"tax_amount"`
	Items          []CreateInvoiceItemReq `json:"items" binding:"required,min=1,dive"`
}

type UpdateInvoiceStatusReq struct {
	Status string `json:"status" binding:"required,oneof=draft finalized cancelled"`
}

type RecordPaymentReq struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Method      string  `json:"method" binding:"required,oneof=cash card online bank_transfer"`
	ReferenceNo *string `json:"reference_no"`
}

type InvoiceItemResponse struct {
	ID              uuid.UUID  `json:"id"`
	BundleID        *uuid.UUID `json:"bundle_id"`
	ItemType        string     `json:"item_type"`
	Description     string     `json:"description"`
	InventoryItemID *uuid.UUID `json:"inventory_item_id"`
	Quantity        int        `json:"quantity"`
	UnitPrice       float64    `json:"unit_price"`
	DiscountAmount  float64    `json:"discount_amount"`
	LineTotal       float64    `json:"line_total"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PaymentResponse struct {
	ID          uuid.UUID `json:"id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	Amount      float64   `json:"amount"`
	Method      string    `json:"method"`
	ReferenceNo *string   `json:"reference_no"`
	PaidAt      time.Time `json:"paid_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type InvoiceResponse struct {
	ID             uuid.UUID             `json:"id"`
	InvoiceNo      string                `json:"invoice_no"`
	PatientID      uuid.UUID             `json:"patient_id"`
	PatientName    string                `json:"patient_name"`
	PatientPhone   string                `json:"patient_phone"`
	VisitID        *uuid.UUID            `json:"visit_id"`
	Status         string                `json:"status"`
	Subtotal       float64               `json:"subtotal"`
	DiscountAmount float64               `json:"discount_amount"`
	TaxAmount      float64               `json:"tax_amount"`
	TotalAmount    float64               `json:"total_amount"`
	PaidAmount     float64               `json:"paid_amount"`
	DueAmount      float64               `json:"due_amount"`
	PaymentStatus  string                `json:"payment_status"`
	CreatedBy      uuid.UUID             `json:"created_by"`
	UpdatedBy      uuid.UUID             `json:"updated_by"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Items          []InvoiceItemResponse `json:"items,omitempty"`
	Payments       []PaymentResponse     `json:"payments,omitempty"`
}

func InvoiceItemRes(it *domain.InvoiceItem) InvoiceItemResponse {
	return InvoiceItemResponse{
		ID:              it.ID,
		BundleID:        it.BundleID,
		ItemType:        string(it.ItemType),
		Description:     it.Description,
		InventoryItemID: it.InventoryItemID,
		Quantity:        it.Quantity,
		UnitPrice:       it.UnitPrice,
		DiscountAmount:  it.DiscountAmount,
		LineTotal:       it.LineTotal,
		CreatedAt:       it.CreatedAt,
	}
}

func PaymentRes(p *domain.Payment) PaymentResponse {
	return PaymentResponse{
		ID:          p.ID,
		InvoiceID:   p.InvoiceID,
		Amount:      p.Amount,
		Method:      string(p.Method),
		ReferenceNo: p.ReferenceNo,
		PaidAt:      p.PaidAt,
		CreatedAt:   p.CreatedAt,
	}
}

func InvoiceRes(d *domain.InvoiceDetails) *InvoiceResponse {
	res := &InvoiceResponse{
		ID:             d.ID,
		InvoiceNo:      d.InvoiceNo,
		PatientID:      d.PatientID,
		PatientName:    d.PatientName,
		PatientPhone:   d.PatientPhone,
		VisitID:        d.VisitID,
		Status:         string(d.Status),
		Subtotal:       d.Subtotal,
		DiscountAmount: d.DiscountAmount,
		TaxAmount:      d.TaxAmount,
		TotalAmount:    d.TotalAmount,
		PaidAmount:     d.PaidAmount,
		DueAmount:      d.DueAmount,
		PaymentStatus:  string(d.PaymentStatus),
		CreatedBy:      d.CreatedBy,
		UpdatedBy:      d.UpdatedBy,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
	for _, it := range d.Items {
		res.Items = append(res.Items, InvoiceItemRes(it))
	}
	for _, p := range d.Payments {
		res.Payments = append(res.Payments, PaymentRes(p))
	}
	return res
}

func InvoiceResList(ds []*domain.InvoiceDetails) []*InvoiceResponse {
	items := make([]*InvoiceResponse, 0, len(ds))
	for _, d := range ds {
		items = append(items, InvoiceRes(d))
	}
	return items
}
