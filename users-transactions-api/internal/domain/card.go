package domain

import "time"

type Card struct {
	Id int32 `json:"id"`
	TenantId int32 `json:"tenant_id"`
	Amount float32 `json:"amount"`
	IdAccount int32 `json:"id_account"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}