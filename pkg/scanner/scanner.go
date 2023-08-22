package scanner

import (
	"database/sql"
	"fmt"
	"os"

	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/sadhasivam/pii-db-scanner/pkg/exporter"
)

const (
	mysqlConnectionFlt    = "%s:%s@tcp(%s:%d)/%s"
	postgresConnectionFlt = "user=%s password=%s host=%s port=%d dbname=%s sslmode=disable"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Scanner interface {
	GetAllSchemas() ([]string, error)
	GetTablesForSchema(schema string) ([]string, error)
	GetTopRecords(schema string, table string) (*sql.Rows, error)
	PrintTable(rows *sql.Rows, exporter exporter.Exporter)
}

func NewScanner(dbType string, db *sql.DB) Scanner {
	if dbType == "mysql" {
		return &MySQLScanner{db: db}
	}
	if dbType == "postgres" {
		return &PostgresScanner{db: db}
	}
	return &DefaultScanner{}
}

func Scan(dbType string, username string, password string, host string, dbport int, dbname string) (Scanner, error) {
	connStr := ""
	if dbType == "mysql" {
		connStr = fmt.Sprintf(mysqlConnectionFlt, username, password, host, dbport, dbname)
	} else if dbType == "postgres" {
		connStr = fmt.Sprintf(postgresConnectionFlt, username, password, host, dbport, dbname)
	}
	fmt.Print(connStr)
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		logger.Error("error obtaining %s database connection host: %s, dbname: %s \n %v ", dbType, host, dbname, err)
		return nil, fmt.Errorf("error obtaining %s database connection host: %s, dbname: %s ", dbType, host, dbname)
	}
	err = db.Ping()
	if err != nil {
		logger.Error("error obtaining %s database connection host: %s, dbname: %s \n %v ", dbType, host, dbname, err)
		return nil, fmt.Errorf("error obtaining %s database connection host: %s, dbname: %s ", dbType, host, dbname)
	}
	return NewScanner(dbType, db), nil
}
