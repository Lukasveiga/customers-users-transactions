package utils

import (
	"github.com/Lukasveiga/customers-users-transactions/internal/ports"
	"github.com/jung-kurt/gofpdf"
)

type GofpdfGenerator struct{}

func NewGofpdfGenerator() *GofpdfGenerator {
	return &GofpdfGenerator{}
}

func (g *GofpdfGenerator) Generate(input ports.PdfGeneratorInputParams) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {
		pdf.SetFont(input.Font, "B", input.FontSize)
		pdf.Cell(0, 10, input.Title)
		pdf.Ln(20)
	})
	pdf.AddPage()

	columnWidths, marginLeft := tableConfig(pdf, 60)

	drawHeaders(pdf, input.Headers, input.Font, input.FontSize, columnWidths, marginLeft)

	// Add table rows
	pdf.SetFont(input.Font, "", input.FontSize)
	for _, data := range input.Data {
		pdf.SetX(marginLeft)
		for j, d := range data {
			pdf.CellFormat(columnWidths[j], 10, d, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		if pdf.GetY() > 260 {
			pdf.AddPage()
			drawHeaders(pdf, input.Headers, input.Font, input.FontSize, columnWidths, marginLeft)
		}
	}

	filePath := "./internal/report-transactions/transactions.pdf"

	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func drawHeaders(pdf *gofpdf.Fpdf, headers []string, font string, fontSize float64, columnWidths []float64,
	marginLeft float64) {
	pdf.SetFont(font, "B", fontSize)
	pdf.SetX(marginLeft)
	for i, header := range headers {
		pdf.CellFormat(columnWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont(font, "", fontSize) // reset font
}

func tableConfig(pdf *gofpdf.Fpdf, cellSize float64) ([]float64, float64) {
	columnWidths := []float64{cellSize, cellSize, cellSize}
	tableWidth := 0.0

	for _, width := range columnWidths {
		tableWidth += width
	}

	pageWidth, _ := pdf.GetPageSize()

	marginLeft := (pageWidth - tableWidth) / 2

	return columnWidths, marginLeft
}
