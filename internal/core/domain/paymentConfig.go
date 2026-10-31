package domain

import "github.com/google/uuid"

// PaymentConfig is what a clinic accepts from its own patients at checkout —
// distinct from the clinic's OpthalmOK subscription billing. BankDetails and
// QRWallet are only meaningful when the matching Methods flag is enabled.
type PaymentConfig struct {
	ClinicID   uuid.UUID
	TaxLabel   string
	TaxRate    float64
	VATNumber  *string
	PANNumber  *string
	Methods    PaymentMethods
	BankDetail *BankDetails
	QRWallet   *QRWallet
}

type PaymentMethods struct {
	Cash         bool
	Card         bool
	BankTransfer bool
	QRWalletFlag bool
}

type BankDetails struct {
	BankName      string
	AccountHolder string
	AccountNumber string
	Branch        *string
	SwiftCode     *string
}

type QRWallet struct {
	Provider   string // eSewa | Khalti | FonePay | other
	MerchantID string
	QRAssetURL *string
}
