package scanner

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sadhasivam/pii-db-scanner/pkg/exporter"
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

func (d *DefaultScanner) getAllSchemas(db *sql.DB, sql string) ([]string, error) {
	rows, err := db.Query(sql)
	if err != nil {
		logger.Error("GetAllSchemas error fetching schema metadata %v ", err)
		return nil, fmt.Errorf("GetAllSchemas error fetching schema metadata %v ", err)
	}
	defer rows.Close()
	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			logger.Error("GetAllSchemas error fetching schema record %v ", err)
			return nil, fmt.Errorf("GetAllSchemas error fetching schema record %v ", err)
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func (d *DefaultScanner) getTablesForSchema(db *sql.DB, sql string) ([]string, error) {
	var tables []string
	rows, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("GetTablesForSchema error fetching tables for SQL: %s", sql)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("GetTablesForSchema error scanning tables for SQL: %s", sql)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (d *DefaultScanner) getTopRecords(db *sql.DB, sql string) (*sql.Rows, error) {
	records, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for sql %s failed", sql)
	}

	_, err = records.Columns()
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for sql %s failed", sql)
	}
	return records, nil
}

func (d *DefaultScanner) PrintTable(records *sql.Rows, exporter exporter.Exporter) {
	cols, err := records.Columns()
	if err != nil {
		log.Fatal(err)
		return
	}
	exporter.PrintTableHeader(cols)

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
		exporter.PrintTableRow(cols, columns)
	}
}
