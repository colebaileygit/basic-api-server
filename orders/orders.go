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
	// body, readErr := ioutil.ReadAll(r.Body)
	// if readErr != nil {
	// 	log.Panic(readErr)
	// }

	// log.Printf("%+v\n", string(body))

	// params := types.PlaceOrderParams{}
	// jsonErr := json.Unmarshal(body, &params)
	// if jsonErr != nil {
	// 	// TODO: Output proper error message and code for incorrect params
	// 	log.Panic(jsonErr)
	// }

	// log.Printf("%+v\n", params)

	var params types.PlaceOrderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	order := types.OrderResponse{
		ID: "1",
		Distance: 50,
		Status: "UNASSIGNED",
	}

	c.JSON(http.StatusOK, order)
}