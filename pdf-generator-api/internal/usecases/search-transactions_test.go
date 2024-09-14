package usecases

import (
	"context"
	"io"
	"testing"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"github.com/stretchr/testify/assert"
)

func TestClientSeachTransactionInfo(t *testing.T) {
	t.Parallel()

	filter := &genproto.Filter{
		TenantId:  1,
		AccountId: 1,
	}

	req := &genproto.SearchTransactionInfoRequest{Filter: filter}
	stream, err := client.SearchTransactionInfo(context.Background(), req)
	assert.NoError(t, err)

	found := 0

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		assert.NoError(t, err)
		assert.Equal(t, filter.GetAccountId(), res.GetTransactionInfo().GetAccountId())

		found += 1
	}

	assert.Equal(t, 2, found)
}
