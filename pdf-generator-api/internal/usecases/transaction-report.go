package usecases

import (
	"fmt"

	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/genproto"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/ports"
	"github.com/Lukasveiga/customers-users-transactions/pdf-generator-api/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionReport struct {
	client       genproto.TransactionInfoServiceClient
	pdfGenerator ports.PdfGenerator
}

func NewTransactionReport(client genproto.TransactionInfoServiceClient, pdfGenerator ports.PdfGenerator) *TransactionReport {
	return &TransactionReport{
		client:       client,
		pdfGenerator: pdfGenerator,
	}
}

type GenerateInputParams struct {
	TenantId  int32
	AccountId int32
}

func (r *TransactionReport) GeneratePdfReport(input GenerateInputParams) (string, error) {

	filter := &genproto.Filter{
		TenantId:  uint32(input.TenantId),
		AccountId: uint32(input.AccountId),
	}

	result, err := SearchTransactionInformation(r.client, filter)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", &shared.EntityNotFoundError{
				Message: err.Error(),
			}
		}
		return "", err
	}

	data := convertData(result)

	inputPdf := ports.PdfGeneratorInputParams{
		Title:    fmt.Sprintf("Account %d Transactions Information", input.AccountId),
		Font:     "Arial",
		FontSize: 12,
		Headers:  []string{"Account", "Kind", "Value"},
		Data:     data,
	}

	path, err := r.pdfGenerator.Generate(inputPdf)

	if err != nil {
		return "", err
	}

	return path, nil
}

func convertData(data []*genproto.TransactionInfo) [][]string {
	table := make([][]string, 0)

	for _, d := range data {
		accountId := fmt.Sprintf("%d", d.AccountId)
		value := fmt.Sprintf("%.2f", d.Value)
		table = append(table, []string{accountId, d.Kind, value})
	}

	return table
}
