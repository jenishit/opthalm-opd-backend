package dto

import (
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type UpsertPosConfigReq struct {
	BarcodeSymbology       string   `json:"barcode_symbology" binding:"required,oneof=code39 code128 ean13"`
	AutoGenerateBarcode    bool     `json:"auto_generate_barcode"`
	BarcodeScanningEnabled bool     `json:"barcode_scanning_enabled"`
	AutoPrintReceipt       bool     `json:"auto_print_receipt"`
	RequirePinForRefunds   bool     `json:"require_pin_for_refunds"`
	ReceiptPrinter         *string  `json:"receipt_printer"`
	TillFloatAmount        *float64 `json:"till_float_amount" binding:"omitempty,min=0"`
}

type PosConfigResponse struct {
	ClinicID               uuid.UUID `json:"clinic_id"`
	BarcodeSymbology       string    `json:"barcode_symbology"`
	AutoGenerateBarcode    bool      `json:"auto_generate_barcode"`
	BarcodeScanningEnabled bool      `json:"barcode_scanning_enabled"`
	AutoPrintReceipt       bool      `json:"auto_print_receipt"`
	RequirePinForRefunds   bool      `json:"require_pin_for_refunds"`
	ReceiptPrinter         *string   `json:"receipt_printer,omitempty"`
	TillFloatAmount        *float64  `json:"till_float_amount,omitempty"`
}

func PosConfigRes(c *domain.PosConfig) *PosConfigResponse {
	if c == nil {
		return nil
	}

	return &PosConfigResponse{
		ClinicID:               c.ClinicID,
		BarcodeSymbology:       c.BarcodeSymbology,
		AutoGenerateBarcode:    c.AutoGenerateBarcode,
		BarcodeScanningEnabled: c.BarcodeScanningEnabled,
		AutoPrintReceipt:       c.AutoPrintReceipt,
		RequirePinForRefunds:   c.RequirePinForRefunds,
		ReceiptPrinter:         c.ReceiptPrinter,
		TillFloatAmount:        c.TillFloatAmount,
	}
}
