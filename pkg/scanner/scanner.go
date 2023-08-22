package scanner

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	mysqlConnectionFlt    = "%s:%s@tcp(%s:%d)/%s"
	postgresConnectionFlt = "user=%s password=%s host=%s port=%d dbname=%s sslmode=disable"
)

type Scanner interface {
	GetAllSchemas() ([]string, error)
	GetTablesForSchema(schema string) ([]string, error)
	GetTopRecords(schema string, table string) (*sql.Rows, error)
	PrintTable(*sql.Rows)
}

func NewScanner(dbType string, db *sql.DB) Scanner {
	if dbType == "mysql" {
		fmt.Println("print mysql")
		return &MySQLScanner{db: db}
	}
	if dbType == "postgres" {
		return &PostgresScanner{db: db}
	}
	return nil
}

func Scan(dbType string, username string, password string, host string, dbport int, dbname string) (Scanner, error) {
	connStr := ""
	if dbType == "mysql" {
		connStr = fmt.Sprintf(mysqlConnectionFlt, username, password, host, dbport, dbname)
		fmt.Println(connStr)
	} else if dbType == "postgres" {
		fmt.Println("print postgres")
		connStr = fmt.Sprintf(postgresConnectionFlt, username, password, host, dbport, dbname)
	}
	fmt.Print(connStr)
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		panic(err)
		return nil, fmt.Errorf("error obtaining %s database connection host: %s, dbname: %s ", dbType, host, dbname)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return NewScanner(dbType, db), nil
}
