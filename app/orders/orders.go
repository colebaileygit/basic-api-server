package orders

import (
	// "log"
	"log"
	"net/http"
	"strconv"

	// "io/ioutil"
	// "encoding/json"

	"github.com/gin-gonic/gin"

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

	// TODO: Execute DB transaction

	order := types.OrderResponse{
		ID:       "1",
		Distance: distance,
		Status:   "UNASSIGNED",
	}

	c.JSON(http.StatusOK, order)
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

// TODO: Assign order

// TODO: Fetch orders
