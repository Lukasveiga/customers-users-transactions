package domain

import (
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type Transaction struct {
	Id        int32     `json:"id"`
	Kind      string    `json:"kind"`
	Value     float32   `json:"value"`
	CardId    int32     `json:"card_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (t *Transaction) Create() (*Transaction, error) {
	err := t.validate()

	if err != nil {
		return nil, err
	}

	t.CreatedAt = time.Now().UTC()
	return t, nil
}

func (t *Transaction) validate() error {
	valErr := &shared.ValidationError{
		Errors: make(map[string]string),
	}

	if len(t.Kind) == 0 {
		valErr.AddError("kind", "cannot be empty")
	}

	if t.Value < 0 {
		valErr.AddError("value", "must be greater than zero (0)")
	}

	if valErr.HasErrors() {
		return valErr
	}

	return nil
}
