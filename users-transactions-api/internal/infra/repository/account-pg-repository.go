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

	err := row.Scan(&savedAccount.Id, &savedAccount.TenantId, &savedAccount.Number, &savedAccount.Status,
		&savedAccount.CreatedAt, &savedAccount.UpdatedAt, &savedAccount.DeletedAt)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "Create"),
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
	err := row.Scan(&account.Id, &account.TenantId, &account.Number, &account.Status,
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

func (ar *PgAccountRepository) FindByNumber(tenantId int32, number string) (*domain.Account, error) {
	var account domain.Account

	query := "SELECT * FROM accounts WHERE tenant_id = $1 AND number = $2"

	row := ar.db.QueryRow(query, tenantId, number)
	err := row.Scan(&account.Id, &account.TenantId, &account.Number, &account.Status,
		&account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil

	case err != nil:
		slog.Error(
			"postgre account repository",
			slog.String("method", "FindByNumber"),
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

		if err := rows.Scan(&account.Id, &account.TenantId, &account.Number, &account.Status,
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

func (ar *PgAccountRepository) Update(id int32, account *domain.Account) (*domain.Account, error) {
	query := `UPDATE accounts SET tenant_id = $1, number = $2, status = $3, created_at = $4, 
			updated_at = $5, deleted_at = $6 WHERE id = $7 RETURNING *`

	row := ar.db.QueryRow(query, account.TenantId, account.Number, account.Status, account.CreatedAt,
		account.UpdatedAt, account.DeletedAt, id)

	var updatedAccount domain.Account

	err := row.Scan(
		&updatedAccount.Id,
		&updatedAccount.TenantId,
		&updatedAccount.Number,
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

func (ar *PgAccountRepository) Delete(tenantId int32, id int32) error {
	query := "UPDATE accounts SET status = 'inactive', deleted_at = NOW() WHERE id = $1 AND tenant_id = $2"

	result, err := ar.db.Exec(query, id, tenantId)

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "Delete"),
			slog.String("error", err.Error()),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		slog.Error(
			"postgre account repository",
			slog.String("method", "Delete.rowsAffected"),
			slog.String("error", err.Error()),
		)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
