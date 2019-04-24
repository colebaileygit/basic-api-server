package main

import (
	"fmt"
	"testing"
	"log"
	"strings"
	"encoding/json"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	"github.com/colebaileygit/basic-api-server/database"
	"github.com/colebaileygit/basic-api-server/types"
)

func setupTestDB(t *testing.T) func(t *testing.T) {
	if testing.Short() {
        t.Skip("Skipping integration test in short mode.")
	}
	
	runDatabaseMigrations()
	return func(t *testing.T) {
		m := database.InitMigrator()

		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while syncing the database: %v", err)
		}
	}
}

func TestPlaceOrder(t *testing.T) {
	testCases := []struct{
		payload string
	}{
		{`{"origin":["22.3082", "114.1765"], "destination": ["22.28653", "114.1789"]}`},
		{`{"origin":["22.25091", "114.0123"], "destination": ["22.2768", "114.1982"]}`},
		{`{"origin":["22.22685", "114.1982"], "destination": ["22.21988", "114.2383"]}`},
		{`{"origin":["22.29084", "114.16541"], "destination": ["22.27653", "114.2109"]}`},
		{`{"origin":["22.28552", "114.15769"], "destination": ["22.3082", "114.1765"]}`},
		{`{"origin":["22.27653", "114.2109"], "destination": ["22.318", "114.183"]}`},
		{`{"origin":["22.265748", "114.1902"], "destination": ["22.29084", "114.16541"]}`},
		{`{"origin":["22.21988", "114.2883"], "destination": ["22.265748", "114.1902"]}`},
		{`{"origin":["22.318", "114.183"], "destination": ["22.25091", "114.0123"]}`},
	}

	testDBTeardown := setupTestDB(t)
	defer testDBTeardown(t)

	router := Routes()

	ids := map[int64]bool {}

	for _, testCase := range testCases {
		t.Run(testCase.payload, func(t *testing.T) {
			response := executeNewRequest(router, "POST", "/orders", testCase.payload)
			checkResponseCode(t, 200, response.Code)

			order := types.OrderResponse{}
			jsonErr := json.Unmarshal(response.Body.Bytes(), &order)
			if jsonErr != nil {
				t.Error("Expected body to be valid order\n")
			}

			if ids[order.ID] {
				t.Errorf("Generated ID %d was duplicated.", order.ID)
			}
			if order.Distance <= 0 {
				t.Errorf("Distance %d must be greater than 0.", order.Distance)
			}
			if order.Status != "UNASSIGNED" {
				t.Errorf("Order status %s must be UNASSIGNED", order.Status)
			}

			ids[order.ID] = true
		})
	}
}

func TestTakeAndFetchOrders(t *testing.T) {
	testCases := []struct{
		page int
		limit int
	}{
		{0, 2},
		{0, 4},
		{1, 2},
		{1, 4},
		{1, 5},
		{2, 4},
		{2, 2},
		{0, 20},
		{0, 1},
		{1, 20},
	}

	testDBTeardown := setupTestDB(t)
	defer testDBTeardown(t)

	router := Routes()

	orderIds := populateOrderTable(t, router)

	// Take a random subset of orders
	takeOrderIds := []int64{
		orderIds[0], orderIds[2], orderIds[3], orderIds[5],
	}

	for _, orderId := range takeOrderIds {
		url := fmt.Sprintf("/orders/%d", orderId)
		response := executeNewRequest(router, "PATCH", url, `{"status": "TAKEN"}`)
		checkResponseCode(t, 200, response.Code)
		checkBody(t, types.TakeOrderResponse{}, response)
	}

	// Test number of retrieved orders is correct and individual order status is correct
	orderStatuses := map[int64]string {}
	for _, id := range orderIds {
		orderStatuses[id] = "UNASSIGNED"
	}
	for _, id := range takeOrderIds {
		orderStatuses[id] = "TAKEN"
	}

	for _, testCase := range testCases {
		expectedNumResults := len(orderIds) - testCase.page * testCase.limit
		if expectedNumResults > testCase.limit {
			expectedNumResults = testCase.limit
		} else if expectedNumResults < 0 {
			expectedNumResults = 0
		}

		url := fmt.Sprintf("/orders?page=%d&limit=%d", testCase.page, testCase.limit)
		response := executeNewRequest(router, "GET", url, "")
		checkResponseCode(t, 200, response.Code)

		orders := types.FetchOrdersResponse{}
		jsonErr := json.Unmarshal(response.Body.Bytes(), &orders)
		if jsonErr != nil {
			t.Error("Expected body to be valid order\n")
		}

		if len(orders) != expectedNumResults {
			t.Errorf("Expected %d results but got %d: %+v\n", expectedNumResults, len(orders), orders)
		}

		for _, order := range orders {
			if order.Status != orderStatuses[order.ID] {
				t.Errorf("Expected %s but retrieved status %s for id %d\n", orderStatuses[order.ID], order.Status, order.ID)
			}
		}
	}


}

func populateOrderTable(t *testing.T, router *gin.Engine) []int64 {
	payloads := []string{
		`{"origin":["22.3082", "114.1765"], "destination": ["22.28653", "114.1789"]}`,
		`{"origin":["22.25091", "114.0123"], "destination": ["22.2768", "114.1982"]}`,
		`{"origin":["22.22685", "114.1982"], "destination": ["22.21988", "114.2383"]}`,
		`{"origin":["22.29084", "114.16541"], "destination": ["22.27653", "114.2109"]}`,
		`{"origin":["22.28552", "114.15769"], "destination": ["22.3082", "114.1765"]}`,
		`{"origin":["22.27653", "114.2109"], "destination": ["22.318", "114.183"]}`,
		`{"origin":["22.265748", "114.1902"], "destination": ["22.29084", "114.16541"]}`,
		`{"origin":["22.21988", "114.2883"], "destination": ["22.265748", "114.1902"]}`,
		`{"origin":["22.318", "114.183"], "destination": ["22.25091", "114.0123"]}`,
	}

	ids := []int64{}

	for _, payload := range payloads {
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(payload))
		response := executeRequest(router, req)
		checkResponseCode(t, 200, response.Code)

		order := types.OrderResponse{}
		jsonErr := json.Unmarshal(response.Body.Bytes(), &order)
		if jsonErr != nil {
			t.Error("Expected body to be valid order\n")
		}

		ids = append(ids, order.ID)
	}

	return ids
}