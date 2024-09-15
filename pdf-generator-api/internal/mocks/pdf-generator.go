package mocks

import (
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/ports"
	"github.com/stretchr/testify/mock"
)

type MockPdfGenerator struct {
	mock.Mock
}

func (mock *MockPdfGenerator) Generate(input ports.PdfGeneratorInputParams) (string, error) {
	args := mock.Called()
	result := args.Get(0)

	return result.(string), args.Error(1)
}
