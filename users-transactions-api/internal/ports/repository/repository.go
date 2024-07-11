package port

import "github.com/Lukasveiga/customers-users-Transaction/internal/domain"

type AccountRepository interface {
	Create(account *domain.Account) (*domain.Account, error)
	FindAll(tenantId int32) ([]domain.Account, error)
	FindByNumber(tenantId int32, number string) (*domain.Account, error)
	Update(id int32, account *domain.Account) (*domain.Account, error)
	Delete(id int32) error
}

type CardRepository interface {
	Create(card *domain.Card) (*domain.Card, error)
	FindAllByAccountId(tenantId int32, accountId int32) ([]domain.Card, error)
	FindById(tenantId int32, id int32) (*domain.Card, error)
	Update(id int32, card *domain.Card) (*domain.Card, error)
	Delete(id int32) error
}

type TransactionRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
	FindAllByCardId(tenantId int32, cardId int32) ([]domain.Transaction, error)
	FindById(tenantId int32, id int32) (*domain.Transaction, error)
	Update(id int32, transaction *domain.Transaction) (*domain.Transaction, error)
	Delete(id int32) error
}

