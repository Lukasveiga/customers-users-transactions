package port

import "github.com/Lukasveiga/customers-users-transaction/internal/domain"

type AccountRepository interface {
	Save(account *domain.Account) (*domain.Account, error)
	FindAll(tenantId int32) ([]*domain.Account, error)
	FindById(tenantId int32, id int32) (*domain.Account, error)
	Update(account *domain.Account) (*domain.Account, error)
}

type CardRepository interface {
	Save(card *domain.Card) (*domain.Card, error)
	FindById(accountId int32, id int32) (*domain.Card, error)
	FindAllByAccountId(accountId int32) ([]*domain.Card, error)
	Update(id int32, card *domain.Card) (*domain.Card, error)
	//Delete(id int32) error
}

type TransactionRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
	FindAllByCardId(tenantId int32, cardId int32) ([]domain.Transaction, error)
	FindById(tenantId int32, id int32) (*domain.Transaction, error)
	Update(id int32, transaction *domain.Transaction) (*domain.Transaction, error)
	Delete(id int32) error
}

type TenantRepository interface {
	FindById(id int32) (*domain.Tenant, error)
}
