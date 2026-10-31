package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jenishit/opthalm-opd-backend/internal/adapter/storage/postgres"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
)

type IntegrationRepository struct {
	DB *postgres.DB
}

func NewIntegrationRepository(db *postgres.DB) *IntegrationRepository {
	return &IntegrationRepository{DB: db}
}

// ListIntegrations always returns one entry per domain.KnownIntegrationProviders,
// even for a clinic that has never connected anything — providers without a
// row come back as Connected: false. It never selects access_token/
// refresh_token, so there's no path for them to reach an API response.
func (r *IntegrationRepository) ListIntegrations(ctx context.Context, clinicID uuid.UUID) ([]*domain.Integration, error) {
	query, args, err := sq.
		Select("id", "provider", "connected", "connected_at").
		From("clinic_integrations").
		Where(sq.Eq{"clinic_id": clinicID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Integration.List query: %w", err)
	}
	defer rows.Close()

	byProvider := make(map[string]*domain.Integration)

	for rows.Next() {
		var i domain.Integration
		i.ClinicID = clinicID

		if err := rows.Scan(&i.ID, &i.Provider, &i.Connected, &i.ConnectedAt); err != nil {
			return nil, fmt.Errorf("Integration.List scan: %w", err)
		}

		byProvider[i.Provider] = &i
	}

	integrations := make([]*domain.Integration, 0, len(domain.KnownIntegrationProviders))
	for _, provider := range domain.KnownIntegrationProviders {
		if i, ok := byProvider[provider]; ok {
			integrations = append(integrations, i)
			continue
		}
		integrations = append(integrations, &domain.Integration{ClinicID: clinicID, Provider: provider})
	}

	return integrations, nil
}

// ConnectIntegration is a stub: it marks the integration connected
// immediately rather than running a real OAuth exchange, since no provider
// credentials are configured anywhere in this codebase yet. Swap this for a
// real redirect+callback flow once they are.
func (r *IntegrationRepository) ConnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) (*domain.Integration, error) {
	query, args, err := sq.
		Insert("clinic_integrations").
		Columns("clinic_id", "provider", "connected", "connected_at").
		Values(clinicID, provider, true, sq.Expr("NOW()")).
		Suffix(`ON CONFLICT (clinic_id, provider) DO UPDATE SET
			connected = TRUE,
			connected_at = NOW(),
			updated_at = NOW()
		RETURNING id, provider, connected, connected_at`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Integration.Connect build: %w", err)
	}

	i := &domain.Integration{ClinicID: clinicID}
	if err := r.DB.QueryRow(ctx, query, args...).Scan(&i.ID, &i.Provider, &i.Connected, &i.ConnectedAt); err != nil {
		return nil, fmt.Errorf("Integration.Connect scan: %w", err)
	}

	return i, nil
}

func (r *IntegrationRepository) DisconnectIntegration(ctx context.Context, clinicID uuid.UUID, provider string) error {
	query, args, err := sq.
		Delete("clinic_integrations").
		Where(sq.Eq{"clinic_id": clinicID, "provider": provider}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("Integration.Disconnect build: %w", err)
	}

	if _, err := r.DB.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Integration.Disconnect exec: %w", err)
	}

	return nil
}
