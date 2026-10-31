package domain

import "github.com/google/uuid"

type PosConfig struct {
	ClinicID               uuid.UUID
	BarcodeSymbology       string // code39 | code128 | ean13
	AutoGenerateBarcode    bool
	BarcodeScanningEnabled bool
	AutoPrintReceipt       bool
	RequirePinForRefunds   bool
	ReceiptPrinter         *string
	TillFloatAmount        *float64
}
