package exporter

type Exporter interface {
	PrintTableHeader(cols []string)
	PrintTableRow(cols []string, columns []interface{})
}

func NewExporter(t string) Exporter {
	return &ConsoleExporter{}
}
