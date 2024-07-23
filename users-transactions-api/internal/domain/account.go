package domain

import (
	"time"
)

type Account struct {
	Id        int32     `json:"id"`
	TenantId  int32     `json:"tenant_id"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
