package infra

import (
	"database/sql"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
)

type PgTenantRepository struct {
	db *sql.DB
}

func NewPgTenantRepository(db *sql.DB) *PgTenantRepository {
	return &PgTenantRepository{
		db: db,
	}
}

func (tr *PgTenantRepository) FindById(id int32) (*domain.Tenant, error) {
	var tenant domain.Tenant

	query := "SELECT * FROM tenants WHERE id = $1"

	row := tr.db.QueryRow(query, id)
	err := row.Scan(&tenant.Id, &tenant.Name)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil

	case err != nil:
		slog.Error(
			"postgre tenant repository",
			slog.String("method", "FindById"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &tenant, nil
}
