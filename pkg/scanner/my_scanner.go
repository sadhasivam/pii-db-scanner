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
	schemas, err := my.getAllSchemas(my.db, "SHOW DATABASES")
	if err == nil {
		// Filtering out system schemas
		var filteredSchemas []string
		for _, schema := range schemas {
			if !isSystemSchema(schema) {
				filteredSchemas = append(filteredSchemas, schema)
			}
		}
		return filteredSchemas, nil
	}
	return schemas, err
}

func (my *MySQLScanner) GetTablesForSchema(schema string) ([]string, error) {
	query := fmt.Sprintf("SHOW TABLES FROM `%s`", schema)
	return my.getTablesForSchema(my.db, query)
}

func (my *MySQLScanner) GetTopRecords(schema string, table string) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 25", schema, table)
	return my.getTopRecords(my.db, sql)
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
