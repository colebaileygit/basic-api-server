package types

type ErrorResponse struct {
	Description string `json:"error"`
}

type OrderResponse struct {
	ID       int64  `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

type PlaceOrderParams struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type TakeOrderParams struct {
	Status string `json:"status"`
}

type TakeOrderResponse struct {
	Status string `json:"status"`
}

type FetchOrdersResponse []OrderResponse
