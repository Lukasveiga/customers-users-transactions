package domain

import (
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type Card struct {
	Id        int32     `json:"id"`
	Amount    float32   `json:"amount"`
	AccountId int32     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (c *Card) Create() *Card {
	c.Amount = float32(0)
	c.CreatedAt = time.Now().UTC()
	return c
}

func (c *Card) AddAmount(value float32) error {
	if value < 0 {
		return &shared.ValidationError{
			Errors: map[string]string{"value": "Amount value must be positive"},
		}
	}

	c.UpdatedAt = time.Now().UTC()
	c.Amount += value
	return nil
}
