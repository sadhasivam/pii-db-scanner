package scanner

import (
	"database/sql"
	"fmt"
)

type PostgresScanner struct {
	DefaultScanner
	db *sql.DB
}

func (pg *PostgresScanner) GetAllSchemas() ([]string, error) {
	// Get all schemas
	fmt.Println("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'pg_toast', 'information_schema') ")
	rows, err := pg.db.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'pg_toast', 'information_schema')")
	if err != nil {
		return nil, fmt.Errorf("GetAllSchemas error fetching schema metadata %v ", err)
	}
	defer rows.Close()
	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, fmt.Errorf("GetAllSchemas error fetching schema record %v ", err)
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func (pg *PostgresScanner) GetTablesForSchema(schema string) ([]string, error) {
	var tables []string

	rows, err := pg.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schema)
	if err != nil {
		return nil, fmt.Errorf("GetTablesForSchema error fetching tables from schema: %s", schema)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("GetTablesForSchema error scanning tables from schema: %s", schema)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (pg *PostgresScanner) GetTopRecords(schema string, table string) (*sql.Rows, error) {
	queryString := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 5", schema, table)

	records, err := pg.db.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for %s.%s not failed", schema, table)
	}

	_, err = records.Columns()
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for %s.%s not failed", schema, table)
	}
	return records, nil
}
