package domain

import "time"

type Transaction struct {
	Id int32 `json:"id"`
	TenantId int32 `json:"tenant_id"`
	Kind string `json:"kind"`
	Value float32 `json:"value"`
	IdCard int32 `json:"id_card"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}