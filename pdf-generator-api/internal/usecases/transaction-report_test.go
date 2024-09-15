package usecases

import (
	"errors"
	"testing"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/mocks"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/shared"
	"github.com/stretchr/testify/assert"
)

func TestTransactionReport(t *testing.T) {
	t.Parallel()

	mockPdfGenerator := new(mocks.MockPdfGenerator)
	sut := NewTransactionReport(client, mockPdfGenerator)

	t.Run("Error to find entities", func(t *testing.T) {
		input := GenerateInputParams{
			TenantId:  1,
			AccountId: 99,
		}

		result, err := sut.GeneratePdfReport(input)

		_, ok := err.(*shared.EntityNotFoundError)

		assert.Empty(t, result)
		assert.True(t, ok)
	})

	t.Run("Error to generate pdf", func(t *testing.T) {
		expectedErr := errors.New("pdf generator error")

		mockPdfGenerator.On("Generate").Return("", expectedErr)
		defer mockPdfGenerator.On("Generate").Unset()

		input := GenerateInputParams{
			TenantId:  1,
			AccountId: 1,
		}

		result, err := sut.GeneratePdfReport(input)

		assert.Empty(t, result)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success", func(t *testing.T) {
		path := "/reports/transactions.pdf"

		mockPdfGenerator.On("Generate").Return(path, nil)
		defer mockPdfGenerator.On("Generate").Unset()

		input := GenerateInputParams{
			TenantId:  1,
			AccountId: 1,
		}

		result, err := sut.GeneratePdfReport(input)

		assert.NoError(t, err)
		assert.Equal(t, path, result)
	})
}
