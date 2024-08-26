package infra

import (
	"database/sql"
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
)

type PgCardRepository struct {
	db *sql.DB
}

func NewPgCardRepository(db *sql.DB) *PgCardRepository {
	return &PgCardRepository{
		db,
	}
}

func (cr *PgCardRepository) Save(card *domain.Card) (*domain.Card, error) {
	var savedCard domain.Card

	query := "INSERT INTO cards (amount, account_id, created_at, updated_at, deleted_at) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING *"

	row := cr.db.QueryRow(query, card.Amount, card.AccountId, card.CreatedAt, card.UpdatedAt,
		card.DeletedAt)

	err := row.Scan(&savedCard.Id, &savedCard.AccountId, &savedCard.Amount, &savedCard.CreatedAt,
		&savedCard.UpdatedAt, &savedCard.DeletedAt)

	if err != nil {
		slog.Error(
			"postgre card repository",
			slog.String("method", "Save"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &savedCard, nil
}

func (cr *PgCardRepository) FindById(accountId int32, id int32) (*domain.Card, error) {
	var card domain.Card

	query := "SELECT * FROM cards WHERE account_id = $1 AND id = $2"

	row := cr.db.QueryRow(query, accountId, id)
	err := row.Scan(&card.Id, &card.Amount, &card.AccountId, &card.CreatedAt,
		&card.UpdatedAt, &card.DeletedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil

	case err != nil:
		slog.Error(
			"postgre card repository",
			slog.String("method", "FindById"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &card, nil
}

func (cr *PgCardRepository) FindAllByAccountId(accountId int32) ([]*domain.Card, error) {
	var cards []*domain.Card

	query := "SELECT * FROM cards WHERE account_id = $1"

	rows, err := cr.db.Query(query, accountId)

	if err != nil {
		slog.Error(
			"postgre card repository",
			slog.String("method", "FindAllByAccountId"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var card domain.Card

		if err := rows.Scan(&card.Id, &card.Amount, &card.AccountId, &card.CreatedAt,
			&card.UpdatedAt, &card.DeletedAt); err != nil {
			slog.Error(
				"postgre card repository",
				slog.String("method", "FindAllByAccountId"),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return cards, nil
	}

	return cards, nil
}

func (cr *PgCardRepository) Update(id int32, card *domain.Card) (*domain.Card, error) {
	query := "UPDATE cards SET amount = $1, account_id = $2, created_at = $3, updated_at = $4, " +
		"deleted_at = $5 WHERE id = $6 RETURNING *"

	row := cr.db.QueryRow(query, &card.Amount, &card.AccountId, &card.CreatedAt,
		&card.UpdatedAt, &card.DeletedAt, id)

	var updatedCard domain.Card

	err := row.Scan(
		&updatedCard.Id,
		&updatedCard.Amount,
		&updatedCard.AccountId,
		&updatedCard.CreatedAt,
		&updatedCard.UpdatedAt,
		&updatedCard.DeletedAt,
	)

	if err != nil {
		slog.Error(
			"postgre card repository",
			slog.String("method", "Update"),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &updatedCard, nil
}
