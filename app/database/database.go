package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	DBCon *sql.DB
)

func Init() {
	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDRESS"), os.Getenv("MYSQL_DATABASE"))

	var err error
	DBCon, err = sql.Open("mysql", mysqlDsn)
	if err != nil {
		log.Fatalf("Could not connect to the MySQL database: %v", err)
	}

	// Ping database until available or max timeout reached
	MAX_TRIES := 30
	SLEEP_TIME := 1 * time.Second
	for i := 1; i <= MAX_TRIES; i++ {
		if err := DBCon.Ping(); err != nil {
			log.Printf("(%d/%d) Could not ping DB, retrying in %s: %v", i, MAX_TRIES, SLEEP_TIME, err)
			if i == MAX_TRIES {
				log.Fatal("Max retries to connect to DB exceeded.")
			}
			time.Sleep(SLEEP_TIME)
		} else {
			return
		}
	}
}

func InitMigrator() *migrate.Migrate {
	if DBCon == nil {
		Init()
		log.Println("Database initialized")
	}

	driver, err := mysql.WithInstance(DBCon, &mysql.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL migration: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "./migrations"),
		"mysql", driver)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return m
}
