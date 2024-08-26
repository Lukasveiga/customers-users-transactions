package domain

import (
	"time"

	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type Account struct {
	Id        int32     `json:"id"`
	TenantId  int32     `json:"tenant_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (a *Account) Create() *Account {
	a.Status = "active"
	a.CreatedAt = time.Now().UTC()
	return a
}

func (a *Account) Active() {
	a.Status = "active"
	a.UpdatedAt = time.Now().UTC()
}

func (a *Account) Inactive() {
	a.Status = "inactive"
	a.DeletedAt = time.Now().UTC()
}

func (a *Account) Validate() error {
	validationErrors := &shared.ValidationError{
		Errors: make(map[string]string),
	}

	if a.Status != "active" && a.Status != "inactive" {
		validationErrors.AddError("status", "must be active or inactive")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
