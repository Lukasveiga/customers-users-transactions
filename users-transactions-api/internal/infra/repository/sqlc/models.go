// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package infra

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int32        `json:"id"`
	TenantID  int32        `json:"tenant_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Card struct {
	ID        int32        `json:"id"`
	AccountID int32        `json:"account_id"`
	Amount    int64        `json:"amount"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Tenant struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID        int32        `json:"id"`
	CardID    int32        `json:"card_id"`
	Kind      string       `json:"kind"`
	Value     int64        `json:"value"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}