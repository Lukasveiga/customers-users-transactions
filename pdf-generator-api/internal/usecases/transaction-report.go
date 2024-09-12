package usecases

import (
	"fmt"

	"github.com/Lukasveiga/customers-users-transactions/internal/genproto"
	"github.com/Lukasveiga/customers-users-transactions/internal/utils"
)

type TransactionReport struct {
	client genproto.TransactionInfoServiceClient
}

func NewTransactionReport(client genproto.TransactionInfoServiceClient) *TransactionReport {
	return &TransactionReport{
		client: client,
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
		// TODO: Custom Error + Logging
		return "", err
	}

	data := convertData(result)

	inputPdf := utils.PdfGeneratorInputParams{
		Title:    fmt.Sprintf("Account %d Transactions Information", input.AccountId),
		Font:     "Arial",
		FontSize: 12,
		Headers:  []string{"Account", "Kind", "Value"},
		Data:     data,
	}

	path, err := utils.PdfGenerator(inputPdf)

	if err != nil {
		// TODO: Custom Error + Logging
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
