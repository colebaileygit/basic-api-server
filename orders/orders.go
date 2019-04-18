package orders

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", PlaceOrder)
	// router.Patch("/{orderId}", TakeOrder)
	// router.Get("/", FetchOrders)

	return router
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	order := Order{
		ID: "1",
		Distance: 50,
		Status: "UNASSIGNED",
	}

	render.JSON(w, r, order)
}

type Order struct {
	ID string `json:"id"`
	Distance int64 `json:"distance"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Description string `json:"error"`
}