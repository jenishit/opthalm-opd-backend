package dto

import (
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PaymentMethodsReq struct {
	Cash         bool `json:"cash"`
	Card         bool `json:"card"`
	BankTransfer bool `json:"bank_transfer"`
	QRWallet     bool `json:"qr_wallet"`
}

type BankDetailsReq struct {
	BankName      string  `json:"bank_name" binding:"required"`
	AccountHolder string  `json:"account_holder" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required"`
	Branch        *string `json:"branch"`
	SwiftCode     *string `json:"swift_code"`
}

type QRWalletReq struct {
	Provider   string  `json:"provider" binding:"required,oneof=eSewa Khalti FonePay other"`
	MerchantID string  `json:"merchant_id" binding:"required"`
	QRAssetURL *string `json:"qr_asset_url"`
}

type UpsertPaymentConfigReq struct {
	TaxLabel   string            `json:"tax_label" binding:"required"`
	TaxRate    float64           `json:"tax_rate" binding:"min=0"`
	VATNumber  *string           `json:"vat_number"`
	PANNumber  *string           `json:"pan_number"`
	Methods    PaymentMethodsReq `json:"methods"`
	BankDetail *BankDetailsReq   `json:"bank_details"`
	QRWallet   *QRWalletReq      `json:"qr_wallet"`
}

type PaymentMethodsResponse struct {
	Cash         bool `json:"cash"`
	Card         bool `json:"card"`
	BankTransfer bool `json:"bank_transfer"`
	QRWallet     bool `json:"qr_wallet"`
}

type BankDetailsResponse struct {
	BankName      string  `json:"bank_name"`
	AccountHolder string  `json:"account_holder"`
	AccountNumber string  `json:"account_number"`
	Branch        *string `json:"branch,omitempty"`
	SwiftCode     *string `json:"swift_code,omitempty"`
}

type QRWalletResponse struct {
	Provider   string  `json:"provider"`
	MerchantID string  `json:"merchant_id"`
	QRAssetURL *string `json:"qr_asset_url,omitempty"`
}

type PaymentConfigResponse struct {
	ClinicID   uuid.UUID              `json:"clinic_id"`
	TaxLabel   string                 `json:"tax_label"`
	TaxRate    float64                `json:"tax_rate"`
	VATNumber  *string                `json:"vat_number,omitempty"`
	PANNumber  *string                `json:"pan_number,omitempty"`
	Methods    PaymentMethodsResponse `json:"methods"`
	BankDetail *BankDetailsResponse   `json:"bank_details,omitempty"`
	QRWallet   *QRWalletResponse      `json:"qr_wallet,omitempty"`
}

func PaymentConfigRes(c *domain.PaymentConfig) *PaymentConfigResponse {
	if c == nil {
		return nil
	}

	rsp := &PaymentConfigResponse{
		ClinicID:  c.ClinicID,
		TaxLabel:  c.TaxLabel,
		TaxRate:   c.TaxRate,
		VATNumber: c.VATNumber,
		PANNumber: c.PANNumber,
		Methods: PaymentMethodsResponse{
			Cash:         c.Methods.Cash,
			Card:         c.Methods.Card,
			BankTransfer: c.Methods.BankTransfer,
			QRWallet:     c.Methods.QRWalletFlag,
		},
	}

	if c.BankDetail != nil {
		rsp.BankDetail = &BankDetailsResponse{
			BankName:      c.BankDetail.BankName,
			AccountHolder: c.BankDetail.AccountHolder,
			AccountNumber: c.BankDetail.AccountNumber,
			Branch:        c.BankDetail.Branch,
			SwiftCode:     c.BankDetail.SwiftCode,
		}
	}

	if c.QRWallet != nil {
		rsp.QRWallet = &QRWalletResponse{
			Provider:   c.QRWallet.Provider,
			MerchantID: c.QRWallet.MerchantID,
			QRAssetURL: c.QRWallet.QRAssetURL,
		}
	}

	return rsp
}
