package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

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

	// TODO: Update port to use ENV variable, add http timeouts etc.
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	// log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(router.Run())
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
