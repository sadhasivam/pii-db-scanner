package exporter

import "fmt"

type ConsoleExporter struct{}

func (c *ConsoleExporter) PrintTableHeader(cols []string) {
	for _, colName := range cols {
		fmt.Printf("%-20s", colName) // Print column name left-aligned in a width of 20 characters
	}
	fmt.Println() // New line after header

	// Print a separator line
	for range cols {
		fmt.Printf("%-20s", "--------------------") // Using '-' for visualization
	}
	fmt.Println()
}

func (c *ConsoleExporter) PrintTableRow(cols []string, columns []interface{}) {
	for _, col := range columns {
		var displayValue string

		switch v := col.(type) {
		case []uint8:
			displayValue = string(v)
		case nil:
			displayValue = "NULL"
		default:
			displayValue = fmt.Sprintf("%v", v)
		}

		fmt.Printf("%-20s", displayValue) // Print value left-aligned in a width of 20 characters
	}
	fmt.Println()
}
