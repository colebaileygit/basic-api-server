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
			Description: err.Error(),
		})
		return
	}

	if params.Origin == nil || params.Destination == nil || len(params.Origin) != 2 || len(params.Destination) != 2 {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: "JSON payload 'origin' or 'destination' was missing or did not contain exactly 2 values [latitude, longitude]",
		})
		return
	}

	order := types.OrderResponse{
		ID:       "1",
		Distance: 50,
		Status:   "UNASSIGNED",
	}

	c.JSON(http.StatusOK, order)
}
