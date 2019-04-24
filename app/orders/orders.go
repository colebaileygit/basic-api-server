package orders

import (
	// "log"
	"log"
	"net/http"
	"strconv"

	// "io/ioutil"
	// "encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/colebaileygit/basic-api-server/database"
	"github.com/colebaileygit/basic-api-server/types"
)

func PlaceOrder(c *gin.Context) {
	var params types.PlaceOrderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "JSON payload could not be parsed.",
		})
		return
	}

	if !validatePlaceOrderParams(params) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "JSON payload 'origin' or 'destination' was missing or did not contain exactly 2 valid values [latitude, longitude]",
		})
		return
	}

	distance, err := CalculateDistance(params)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "Distance between provided origin and destination could not be calculated.",
		})
		return
	}

	if !validateDatabaseConnection(c) {
		return
	}

	orderStatus := "UNASSIGNED"
	res, err := database.DBCon.Exec("INSERT INTO orders (distance, order_status) VALUES (?, ?)",
		distance, orderStatus)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to save order to database.",
		})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to save order to database.",
		})
		return
	}

	order := types.OrderResponse{
		ID:       id,
		Distance: distance,
		Status:   orderStatus,
	}

	c.JSON(http.StatusOK, order)
}

func validateDatabaseConnection(c *gin.Context) bool {
	if database.DBCon == nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Database connection not established.",
		})
		return false
	}
	return true
}

func validatePlaceOrderParams(params types.PlaceOrderParams) bool {
	return params.Origin != nil &&
		params.Destination != nil &&
		len(params.Origin) == 2 &&
		len(params.Destination) == 2 &&
		validateLat(params.Origin[0]) && validateLng(params.Origin[1]) &&
		validateLat(params.Destination[0]) && validateLng(params.Destination[1])

}

func validateLat(val string) bool {
	latitude, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}

	return latitude >= -90 && latitude <= 90
}

func validateLng(val string) bool {
	longitude, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}

	return longitude >= -180 && longitude <= 180
}

func TakeOrder(c *gin.Context) {
	id := c.Param("id")

	var params types.TakeOrderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "JSON payload could not be parsed.",
		})
		return
	}

	if !validateTakeOrderParams(params) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "JSON payload 'status' value was not valid.",
		})
		return
	}

	if !validateDatabaseConnection(c) {
		return
	}

	res, err := database.DBCon.Exec("UPDATE orders SET order_status=? WHERE id=? AND order_status='UNASSIGNED'",
		params.Status, id)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to update order in database.",
		})
		return
	}

	numRows, err := res.RowsAffected()
	if err != nil || numRows == 0 {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to update order in database.",
		})
		return
	}

	response := types.TakeOrderResponse{
		Status: "SUCCESS",
	}

	c.JSON(http.StatusOK, response)
}

func validateTakeOrderParams(params types.TakeOrderParams) bool {
	return params.Status == "TAKEN"
}

func FetchOrders(c *gin.Context) {
	pageArg := c.DefaultQuery("page", "0")
	limitArg := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageArg, 10, 64)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "Invalid page argument. Please provide an integer value >= 0.",
		})
		return
	}

	limit, err := strconv.ParseInt(limitArg, 10, 64)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "Invalid limit argument. Please provide an integer value >= 1.",
		})
		return
	}

	if !validateDatabaseConnection(c) {
		return
	}

	response := types.FetchOrdersResponse{}

	rows, err := database.DBCon.Query("SELECT id, distance, order_status FROM orders LIMIT ?, ?",
		limit*page, limit)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to fetch orders from database.",
		})
		return
	}

	defer rows.Close()
	for rows.Next() {
		order := &types.OrderResponse{}

		err := rows.Scan(&order.ID, &order.Distance, &order.Status)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{
				Description: "Failed to fetch orders from database.",
			})
			return
		}
		response = append(response, *order)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Description: "Failed to fetch orders from database.",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
