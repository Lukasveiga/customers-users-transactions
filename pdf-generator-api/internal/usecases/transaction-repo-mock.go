package usecases

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"github.com/jinzhu/copier"
)

type TransactionInfoRepository interface {
	Search(ctx context.Context, filter *genproto.Filter,
		found func(transactionInfo *genproto.TransactionInfo) error) error
}

type InMemoryTransactionInfoRepository struct {
	mutex sync.RWMutex
	data  []*genproto.TransactionInfo
}

func NewInMemoryTransactionInfoRepository() *InMemoryTransactionInfoRepository {
	return &InMemoryTransactionInfoRepository{
		data: []*genproto.TransactionInfo{
			{
				AccountId: 1,
				Kind:      "Streaming Z",
				Value:     50,
			},
			{
				AccountId: 1,
				Kind:      "Streaming X",
				Value:     60,
			},
			{
				AccountId: 3,
				Kind:      "Streaming Z",
				Value:     50,
			},
		},
	}
}

func (repo *InMemoryTransactionInfoRepository) Search(ctx context.Context, filter *genproto.Filter,
	found func(transactionInfo *genproto.TransactionInfo) error) error {

	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	for _, transInfo := range repo.data {

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			return errors.New("context is canceled")
		}

		if isQualified(filter, transInfo) {
			other, err := deepCopy(transInfo)

			if err != nil {
				return err
			}

			err = found(other)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *genproto.Filter, transInfo *genproto.TransactionInfo) bool {
	return transInfo.GetAccountId() == filter.GetAccountId()
}

func deepCopy(transInfo *genproto.TransactionInfo) (*genproto.TransactionInfo, error) {
	other := &genproto.TransactionInfo{}
	err := copier.Copy(other, transInfo)

	if err != nil {
		return nil, fmt.Errorf("cannot copy transaction data: %w", err)
	}

	return other, nil
}
