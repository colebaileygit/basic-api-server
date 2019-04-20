package orders

import (
	// "log"
	"net/http"
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
			Description: "JSON payload 'origin' or 'destination' was missing or did not contain exactly 2 values [latitude, longitude]",
		})
		return
	}

	// TODO: Calculate distance

	// TODO: Execute DB transaction

	order := types.OrderResponse{
		ID:       "1",
		Distance: 50,
		Status:   "UNASSIGNED",
	}

	c.JSON(http.StatusOK, order)
}

func validatePlaceOrderParams(params types.PlaceOrderParams) bool {
	return params.Origin != nil &&
		params.Destination != nil &&
		len(params.Origin) == 2 &&
		len(params.Destination) == 2
}

// TODO: Assign order

// TODO: Fetch orders
