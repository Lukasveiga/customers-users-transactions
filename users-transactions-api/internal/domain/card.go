package domain

import (
	"time"
)

type Card struct {
	Id        int32     `json:"id"`
	Amount    float32   `json:"amount"`
	AccountId int32     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (a *Card) Create() *Card {
	a.Amount = float32(0)
	a.CreatedAt = time.Now().UTC()
	return a
}
