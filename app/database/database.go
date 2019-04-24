package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	DBCon *sql.DB
)

func Init() {
	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), "db", os.Getenv("MYSQL_DATABASE"))

	var err error
	DBCon, err = sql.Open("mysql", mysqlDsn)
	if err != nil {
		log.Fatalf("Could not connect to the MySQL database: %v", err)
	}

	if err := DBCon.Ping(); err != nil {
		log.Fatalf("Could not ping DB: %v", err)
	}
}
