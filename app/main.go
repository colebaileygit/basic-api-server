package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/colebaileygit/basic-api-server/orders"
	"github.com/colebaileygit/basic-api-server/types"
)

func Routes() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Logger(),
		globalRecover,
	)
	router.NoRoute(notFound)

	ordersEndpoint := router.Group("/orders")
	{
		ordersEndpoint.POST("", orders.PlaceOrder)
		// router.Patch("/{orderId}", TakeOrder)
		// router.Get("/", FetchOrders)
	}

	return router
}

func main() {
	router := Routes()

	runDatabaseMigrations()

	// TODO: Update port to use ENV variable, add http timeouts etc.
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	log.Fatal(router.Run())
}

func runDatabaseMigrations() {
	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), "db", os.Getenv("MYSQL_DATABASE"))

	db, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		log.Fatalf("Could not connect to the MySQL database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping DB: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL migration: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "./migrations"),
		"mysql", driver)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database: %v", err)
	}

	log.Println("Database migrated")
}

// Handle downstream errors and gracefully return 500 error code
func globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			log.Printf("Panic encountered: %+v\n", rec)
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{
				Description: "internal server error", //rec,
			})
		}
	}(c)
	c.Next()
}

// Handle invalid routes or HTTP methods
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, types.ErrorResponse{
		Description: "route is invalid for given request",
	})
}
