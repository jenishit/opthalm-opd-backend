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

type PosConfigRepository struct {
	DB *postgres.DB
}

func NewPosConfigRepository(db *postgres.DB) *PosConfigRepository {
	return &PosConfigRepository{DB: db}
}

func (r *PosConfigRepository) UpsertPosConfig(ctx context.Context, c *domain.PosConfig) (*domain.PosConfig, error) {
	query, args, err := sq.
		Insert("clinic_pos_config").
		Columns(
			"clinic_id", "barcode_symbology", "auto_generate_barcode", "barcode_scanning_enabled",
			"auto_print_receipt", "require_pin_for_refunds", "receipt_printer", "till_float_amount",
		).
		Values(
			c.ClinicID, c.BarcodeSymbology, c.AutoGenerateBarcode, c.BarcodeScanningEnabled,
			c.AutoPrintReceipt, c.RequirePinForRefunds, c.ReceiptPrinter, c.TillFloatAmount,
		).
		Suffix(`ON CONFLICT (clinic_id) DO UPDATE SET
			barcode_symbology = EXCLUDED.barcode_symbology,
			auto_generate_barcode = EXCLUDED.auto_generate_barcode,
			barcode_scanning_enabled = EXCLUDED.barcode_scanning_enabled,
			auto_print_receipt = EXCLUDED.auto_print_receipt,
			require_pin_for_refunds = EXCLUDED.require_pin_for_refunds,
			receipt_printer = EXCLUDED.receipt_printer,
			till_float_amount = EXCLUDED.till_float_amount,
			updated_at = NOW()
		RETURNING clinic_id, barcode_symbology, auto_generate_barcode, barcode_scanning_enabled,
			auto_print_receipt, require_pin_for_refunds, receipt_printer, till_float_amount`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PosConfig.Upsert build: %w", err)
	}

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&c.ClinicID, &c.BarcodeSymbology, &c.AutoGenerateBarcode, &c.BarcodeScanningEnabled,
		&c.AutoPrintReceipt, &c.RequirePinForRefunds, &c.ReceiptPrinter, &c.TillFloatAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("PosConfig.Upsert scan: %w", err)
	}

	return c, nil
}

func (r *PosConfigRepository) GetPosConfig(ctx context.Context, clinicID uuid.UUID) (*domain.PosConfig, error) {
	query, args, err := sq.
		Select(
			"clinic_id", "barcode_symbology", "auto_generate_barcode", "barcode_scanning_enabled",
			"auto_print_receipt", "require_pin_for_refunds", "receipt_printer", "till_float_amount",
		).
		From("clinic_pos_config").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var c domain.PosConfig

	err = r.DB.QueryRow(ctx, query, args...).Scan(
		&c.ClinicID, &c.BarcodeSymbology, &c.AutoGenerateBarcode, &c.BarcodeScanningEnabled,
		&c.AutoPrintReceipt, &c.RequirePinForRefunds, &c.ReceiptPrinter, &c.TillFloatAmount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.PosConfig{
				ClinicID:               clinicID,
				BarcodeSymbology:       "code128",
				AutoGenerateBarcode:    true,
				BarcodeScanningEnabled: true,
			}, nil
		}
		return nil, fmt.Errorf("PosConfig.Get scan: %w", err)
	}

	return &c, nil
}
