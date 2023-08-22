package scanner

import (
	"database/sql"
	"fmt"
	"log"
)

type DefaultScanner struct {
}

func (d *DefaultScanner) GetAllSchemas() ([]string, error) {
	return nil, fmt.Errorf("GetAllSchemas not implemeted")
}

func (d *DefaultScanner) GetTablesForSchema(schema string) ([]string, error) {
	return nil, fmt.Errorf("GetTablesForSchema for %s not implemeted", schema)
}

func (d *DefaultScanner) GetTopRecords(schema string, table string) (*sql.Rows, error) {
	return nil, fmt.Errorf("GetTopRecords for %s.%s not implemeted", schema, table)
}

func printTableHeader(cols []string) {
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

func printTableRow(cols []string, columns []interface{}) {
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

func (d *DefaultScanner) PrintTable(records *sql.Rows) {
	cols, err := records.Columns()
	if err != nil {
		log.Fatal(err)
		return
	}

	printTableHeader(cols)

	for records.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := records.Scan(columnPointers...); err != nil {
			log.Fatal(err)
			continue
		}
		printTableRow(cols, columns)
	}
}
