package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type PaymentConfigRepository struct {
	DB *postgres.DB
}

func NewPaymentConfigRepository(db *postgres.DB) *PaymentConfigRepository {
	return &PaymentConfigRepository{DB: db}
}

func (r *PaymentConfigRepository) UpsertPaymentConfig(ctx context.Context, c *domain.PaymentConfig) (*domain.PaymentConfig, error) {
	var bankName, accountHolder, accountNumber, branch, swiftCode *string
	if c.BankDetail != nil {
		bankName, accountHolder, accountNumber = &c.BankDetail.BankName, &c.BankDetail.AccountHolder, &c.BankDetail.AccountNumber
		branch, swiftCode = c.BankDetail.Branch, c.BankDetail.SwiftCode
	}

	var qrProvider, qrMerchantID *string
	if c.QRWallet != nil {
		qrProvider, qrMerchantID = &c.QRWallet.Provider, &c.QRWallet.MerchantID
	}

	query, args, err := sq.
		Insert("clinic_payment_config").
		Columns(
			"clinic_id", "tax_label", "tax_rate", "vat_number", "pan_number",
			"method_cash", "method_card", "method_bank_transfer", "method_qr_wallet",
			"bank_name", "bank_account_holder", "bank_account_number", "bank_branch", "bank_swift_code",
			"qr_provider", "qr_merchant_id",
		).
		Values(
			c.ClinicID, c.TaxLabel, c.TaxRate, c.VATNumber, c.PANNumber,
			c.Methods.Cash, c.Methods.Card, c.Methods.BankTransfer, c.Methods.QRWalletFlag,
			bankName, accountHolder, accountNumber, branch, swiftCode,
			qrProvider, qrMerchantID,
		).
		Suffix(`ON CONFLICT (clinic_id) DO UPDATE SET
			tax_label = EXCLUDED.tax_label,
			tax_rate = EXCLUDED.tax_rate,
			vat_number = EXCLUDED.vat_number,
			pan_number = EXCLUDED.pan_number,
			method_cash = EXCLUDED.method_cash,
			method_card = EXCLUDED.method_card,
			method_bank_transfer = EXCLUDED.method_bank_transfer,
			method_qr_wallet = EXCLUDED.method_qr_wallet,
			bank_name = EXCLUDED.bank_name,
			bank_account_holder = EXCLUDED.bank_account_holder,
			bank_account_number = EXCLUDED.bank_account_number,
			bank_branch = EXCLUDED.bank_branch,
			bank_swift_code = EXCLUDED.bank_swift_code,
			qr_provider = EXCLUDED.qr_provider,
			qr_merchant_id = EXCLUDED.qr_merchant_id,
			updated_at = NOW()`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PaymentConfig.Upsert build: %w", err)
	}

	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("PaymentConfig.Upsert exec: %w", err)
	}

	return r.GetPaymentConfig(ctx, c.ClinicID)
}

func (r *PaymentConfigRepository) GetPaymentConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PaymentConfig, error) {
	query, args, err := sq.
		Select(
			"clinic_id", "tax_label", "tax_rate", "vat_number", "pan_number",
			"method_cash", "method_card", "method_bank_transfer", "method_qr_wallet",
			"bank_name", "bank_account_holder", "bank_account_number", "bank_branch", "bank_swift_code",
			"qr_provider", "qr_merchant_id", "qr_asset_url",
		).
		From("clinic_payment_config").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	return scanPaymentConfig(r.DB.QueryRow(ctx, query, args...), clinicID)
}

func (r *PaymentConfigRepository) SetQRAssetURL(ctx context.Context, clinicID uuid.UUID, url string) (*domain.PaymentConfig, error) {
	query, args, err := sq.
		Update("clinic_payment_config").
		Set("qr_asset_url", url).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PaymentConfig.SetQRAssetURL build: %w", err)
	}

	tag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PaymentConfig.SetQRAssetURL exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, domain.ErrDataNotFound
	}

	return r.GetPaymentConfig(ctx, clinicID)
}

// defaultPaymentConfig is what a clinic sees before it has ever saved a
// payment configuration — cash-only, zero tax, same as the column defaults.
func defaultPaymentConfig(clinicID uuid.UUID) *domain.PaymentConfig {
	return &domain.PaymentConfig{
		ClinicID: clinicID,
		TaxLabel: "VAT",
		Methods:  domain.PaymentMethods{Cash: true},
	}
}

func scanPaymentConfig(row pgx.Row, clinicID uuid.UUID) (*domain.PaymentConfig, error) {
	var c domain.PaymentConfig
	var bankName, accountHolder, accountNumber, branch, swiftCode *string
	var qrProvider, qrMerchantID, qrAssetURL *string

	err := row.Scan(
		&c.ClinicID, &c.TaxLabel, &c.TaxRate, &c.VATNumber, &c.PANNumber,
		&c.Methods.Cash, &c.Methods.Card, &c.Methods.BankTransfer, &c.Methods.QRWalletFlag,
		&bankName, &accountHolder, &accountNumber, &branch, &swiftCode,
		&qrProvider, &qrMerchantID, &qrAssetURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return defaultPaymentConfig(clinicID), nil
		}
		return nil, fmt.Errorf("PaymentConfig scan: %w", err)
	}

	if bankName != nil {
		c.BankDetail = &domain.BankDetails{
			BankName:      *bankName,
			AccountHolder: *accountHolder,
			AccountNumber: *accountNumber,
			Branch:        branch,
			SwiftCode:     swiftCode,
		}
	}
	if qrProvider != nil {
		c.QRWallet = &domain.QRWallet{
			Provider:   *qrProvider,
			MerchantID: *qrMerchantID,
			QRAssetURL: qrAssetURL,
		}
	}

	return &c, nil
}
