package scanner

import (
	"database/sql"
	"fmt"
)

type MySQLScanner struct {
	DefaultScanner
	db *sql.DB
}

func (my *MySQLScanner) GetAllSchemas() ([]string, error) {
	rows, err := my.db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("GetAllSchemas error fetching show databases %v", err)
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, fmt.Errorf("GetAllSchemas error fetching schema metadata %v", err)
		}
		if !isSystemSchema(schema) {
			schemas = append(schemas, schema)
		}
		fmt.Printf("Skipping analysis for %s as it's a MySQL system schema.\n", schema)
	}
	return schemas, nil
}

func (my *MySQLScanner) GetTablesForSchema(schema string) ([]string, error) {
	var tables []string
	query := fmt.Sprintf("SHOW TABLES FROM `%s`", schema)
	rows, err := my.db.Query(query)
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

func (my *MySQLScanner) GetTopRecords(schema string, table string) (*sql.Rows, error) {
	queryString := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 5", schema, table)

	records, err := my.db.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for %s.%s not failed", schema, table)
	}

	_, err = records.Columns()
	if err != nil {
		return nil, fmt.Errorf("GetTopRecords for %s.%s not failed", schema, table)
	}
	return records, nil
}

func isSystemSchema(schema string) bool {
	systemSchemas := []string{"information_schema", "mysql", "performance_schema", "innodb", "sys"}
	for _, s := range systemSchemas {
		if s == schema {
			return true
		}
	}
	return false
}
