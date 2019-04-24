package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/colebaileygit/basic-api-server/types"
)

// Test that correct routes are setup, including redirects for trailing slashes and 404 for invalid routes.
func TestRoutes(t *testing.T) {
	router := Routes()

	testCases := []struct {
		url     string
		command string
		code    int
		payload string
	}{
		// Valid requests return 500 on unit tests because of missing DB
		{"/orders", 	"POST", 	500, 	`{"origin": ["23", "100"], "destination": ["24", "101"]}`},
		{"/orders/", 	"POST", 	307, 	`{"origin": ["23", "100"], "destination": ["24", "101"]}`},
		{"/orders/0", 	"PATCH", 	500, 	`{"status": "TAKEN"}`},
		{"/orders/0/", 	"PATCH", 	307, 	`{"status": "TAKEN"}`},
		{"/orders", 	"GET", 		500, 	""},
		{"/orders/", 	"GET", 		301, 	""},
		{"/orders/1", 	"POST", 	404, 	""},
		{"/orders", 	"PATCH", 	404, 	""},
		{"/orders/", 	"PATCH", 	404, 	""},
		{"/order", 		"POST", 	404, 	""},
		{"/",			"POST", 	404, 	""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.command+"->"+testCase.url, func(t *testing.T) {
			response := executeNewRequest(router, testCase.command, testCase.url, testCase.payload)

			checkResponseCode(t, testCase.code, response.Code)

			if testCase.code != 307 && testCase.code != 301 {
				success := checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])
				if !success {
					t.Logf("Response body: %s\n", response.Body.String())
				}
			}
		})
	}
}

// Test parsing of user params sent over HTTP
func TestPlaceOrderParams(t *testing.T) {
	router := Routes()

	testCases := []struct {
		description string
		code        int
		payload     string
	}{
		// Valid requests return 500 on unit tests because of missing DB
		{"valid-payload", 500, `{"origin": ["23", "100"], "destination": ["24", "101"]}`},
		{"valid-payload-random-arg", 500, `{"origin": ["23", "100"], "destination": ["24", "101"], "version": 2.0}`},
		{"invalid-payload-integers", 400, `{"origin": [23, 100], "destination": [24, 101]}`},
		{"invalid-payload-floats", 400, `{"origin": [23.0, 100.0], "destination": [24.0, 101.0]}`},
		{"invalid-payload-missing", 400, ""},
		{"invalid-payload-empty", 400, "{}"},
		{"invalid-payload-origin-missing", 400, `{"destination": ["24", "102"]}`},
		{"invalid-payload-origin-empty", 400, `{"origin": [], "destination": ["24", "102"]}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			response := executeNewRequest(router, "POST", "/orders", testCase.payload)

			checkResponseCode(t, testCase.code, response.Code)
			checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])

			switch testCase.code {
			case 200:
				checkBody(t, types.OrderResponse{}, response)
			case 400, 500:
				checkBody(t, types.ErrorResponse{}, response)
			}
		})
	}
}

func TestTakeOrderParams(t *testing.T) {
	router := Routes()

	testCases := []struct {
		description string
		code        int
		payload     string
	}{
		// Valid requests return 500 on unit tests because of missing DB
		{"valid-payload", 500, `{"status": "TAKEN"}`},
		{"valid-payload-random-arg", 500, `{"status": "TAKEN", "version": 2.0}`},
		{"invalid-payload-unassigned", 400, `{"status": "ASSIGNED"}`},
		{"invalid-payload-assigned", 400, `{"status": "UNASSIGNED"}`},
		{"invalid-payload-missing", 400, ""},
		{"invalid-payload-empty", 400, "{}"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			response := executeNewRequest(router, "PATCH", "/orders/1", testCase.payload)

			checkResponseCode(t, testCase.code, response.Code)
			checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])

			switch testCase.code {
			case 200:
				checkBody(t, types.TakeOrderResponse{}, response)
			case 400, 500:
				checkBody(t, types.ErrorResponse{}, response)
			}
		})
	}
}

func TestFetchOrdersParams(t *testing.T) {
	router := Routes()

	testCases := []struct {
		description string
		code        int
		payload     string
	}{
		// Valid requests return 500 on unit tests because of missing DB
		{"valid-payload", 500, `page=0&limit=10`},
		{"valid-payload-random-arg", 500, `page=0&limit=10&version=2.0`},
		{"valid-payload-missing", 500, ""},
		{"invalid-payload-negativepage", 400, `page=-1`},
		{"invalid-payload-zerolimit", 400, `limit=0`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			url := fmt.Sprintf("/orders?%s", testCase.payload)
			response := executeNewRequest(router, "GET", url, "")

			checkResponseCode(t, testCase.code, response.Code)
			checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])

			switch testCase.code {
			case 200:
				checkBody(t, types.FetchOrdersResponse{}, response)
			case 400, 500:
				checkBody(t, types.ErrorResponse{}, response)
			}
		})
	}
}

func executeNewRequest(router *gin.Engine, method string, url string, body string) *httptest.ResponseRecorder {
	requestBody := strings.NewReader(body)
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return executeRequest(router, req)
}

func executeRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, but actual is %d\n", expected, actual)
	}
}

func checkHeader(t *testing.T, expected string, actual []string) bool {
	found := false
	for _, header := range actual {
		if expected == header {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected header %s, but actual is %v\n", expected, actual)
	}
	return found
}

func checkBody(t *testing.T, expected interface{}, request *httptest.ResponseRecorder) {
	jsonErr := json.Unmarshal(request.Body.Bytes(), &expected)
	if jsonErr != nil {
		t.Errorf("Expected body to be valid %T\n", expected)
	}
}
