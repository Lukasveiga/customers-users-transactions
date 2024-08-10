package infra

import (
	"database/sql"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
)

type PgAccountRepository struct {
	db *sql.DB
}

func NewPgAccountRepository(db *sql.DB) *PgAccountRepository {
	return &PgAccountRepository{
		db,
	}
}

func (ar *PgAccountRepository) Create(account *domain.Account) (*domain.Account, error) {
	var savedAccount domain.Account

	query := "INSERT INTO accounts (tenant_id, number, status, created_at, updated_at, deleted_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"

	row := ar.db.QueryRow(query, account.TenantId, account.Number, account.Status, account.CreatedAt,
		account.UpdatedAt, account.DeletedAt)

	err := row.Scan(&savedAccount.Id, &savedAccount.TenantId, &savedAccount.Status, &savedAccount.CreatedAt,
		&savedAccount.UpdatedAt, &savedAccount.DeletedAt)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "create"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &savedAccount, nil
}
