package domain

import (
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
	"github.com/google/uuid"
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

func (a *Account) Validate() error {
	validationErrors := &shared.ValidationError{
		Errors: make(map[string]string),
	}

	_, err := uuid.Parse(a.Number)

	if err != nil {
		validationErrors.AddError("number", "must be a valid uuid")
	}

	if a.Status != "active" && a.Status != "inactive" {
		validationErrors.AddError("status", "must be active or inactive")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
