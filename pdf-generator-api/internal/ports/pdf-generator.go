package ports

type PdfGenerator interface {
	Generate(PdfGeneratorInputParams) (string, error)
}

type PdfGeneratorInputParams struct {
	Title    string
	Font     string
	FontSize float64
	Headers  []string
	Data     [][]string
}
