package dtos

import "github.com/Lukasveiga/customers-users-transaction/internal/domain"

type AccountDto struct {
	Number string `json:"number"`
	Status string `json:"status"`
}

func (ad AccountDto) ToDomain() *domain.Account {
	return &domain.Account{
		Number: ad.Number,
		Status: ad.Status,
	}
}
