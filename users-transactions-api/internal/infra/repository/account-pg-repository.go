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

func (ar *PgAccountRepository) Save(account *domain.Account) (*domain.Account, error) {
	var savedAccount domain.Account

	query := "INSERT INTO accounts (tenant_id, status, created_at, updated_at, deleted_at) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING *"

	row := ar.db.QueryRow(query, account.TenantId, account.Status, account.CreatedAt,
		account.UpdatedAt, account.DeletedAt)

	err := row.Scan(&savedAccount.Id, &savedAccount.TenantId, &savedAccount.Status,
		&savedAccount.CreatedAt, &savedAccount.UpdatedAt, &savedAccount.DeletedAt)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "Save"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &savedAccount, nil
}

func (ar *PgAccountRepository) FindById(tenantId int32, id int32) (*domain.Account, error) {
	var account domain.Account

	query := "SELECT * FROM accounts WHERE tenant_id = $1 AND id = $2"

	row := ar.db.QueryRow(query, tenantId, id)
	err := row.Scan(&account.Id, &account.TenantId, &account.Status,
		&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil

	case err != nil:
		slog.Error(
			"postgre account repository",
			slog.String("method", "FindById"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &account, nil
}

func (ar *PgAccountRepository) FindAll(tenantId int32) ([]*domain.Account, error) {
	var accounts []*domain.Account

	query := "SELECT * FROM accounts WHERE tenant_id = $1"

	rows, err := ar.db.Query(query, tenantId)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "FindAll"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var account domain.Account

		if err := rows.Scan(&account.Id, &account.TenantId, &account.Status,
			&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt); err != nil {
			slog.Error(
				"postgre account repository",
				slog.String("method", "FindAll"),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return accounts, nil
	}

	return accounts, nil
}

func (ar *PgAccountRepository) Update(account *domain.Account) (*domain.Account, error) {
	query := `UPDATE accounts SET tenant_id = $1, status = $2, created_at = $3, 
			updated_at = $4, deleted_at = $5 WHERE id = $6 RETURNING *`

	row := ar.db.QueryRow(query, account.TenantId, account.Status, account.CreatedAt,
		account.UpdatedAt, account.DeletedAt, account.Id)

	var updatedAccount domain.Account

	err := row.Scan(
		&updatedAccount.Id,
		&updatedAccount.TenantId,
		&updatedAccount.Status,
		&updatedAccount.CreatedAt,
		&updatedAccount.UpdatedAt,
		&updatedAccount.DeletedAt,
	)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "Update"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &updatedAccount, nil
}
