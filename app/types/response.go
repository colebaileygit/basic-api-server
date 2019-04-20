package types

type ErrorResponse struct {
	Description string `json:"error"`
}

type OrderResponse struct {
	ID string `json:"id"`
	Distance int64 `json:"distance"`
	Status string `json:"status"`
}

type PlaceOrderParams struct {
	Origin []string `json:"origin"`
	Destination []string `json:"destination"`
}