package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/colebaileygit/basic-api-server/types"
)

func TestRoutes(t *testing.T) {
	router := Routes()

	testCases := []struct {
		url     string
		command string
		code    int
		payload string
	}{
		{"/orders", "POST", 200, `{"origin": ["100", "100"], "destination": ["101", "101"]}`},
		{"/orders/", "POST", 307, `{"origin": ["100", "100"], "destination": ["101", "101"]}`},
		{"/orders/1", "POST", 404, ""},
		// {"/order", 		"POST", 	404, 	""},
		// {"/",			"POST", 	404, 	""},
		// {"/orders", 	"PATCH", 	200, 	""},
		// {"/orders", 	"GET", 		200, 	""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.command+"->"+testCase.url, func(t *testing.T) {
			req, _ := http.NewRequest(testCase.command, testCase.url, strings.NewReader(testCase.payload))
			response := executeRequest(router, req)

			checkResponseCode(t, testCase.code, response.Code)

			if testCase.code != 307 {
				checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])
			}

			t.Logf("Response body: %s\n", response.Body.String())
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	router := Routes()

	testCases := []struct {
		description string
		code        int
		payload     string
	}{
		{"valid-payload", 200, `{"origin": ["100", "100"], "destination": ["101", "101"]}`},
		{"invalid-payload-integers", 400, `{"origin": [100, 100], "destination": [101, 101]}`},
		{"invalid-payload-floats", 400, `{"origin": [100.0, 100.0], "destination": [101.0, 101.0]}`},
		{"invalid-payload-missing", 400, ""},
		{"invalid-payload-empty", 400, "{}"},
		{"invalid-payload-origin-missing", 400, `{"destination": ["101", "102"]}`},
		{"invalid-payload-origin-empty", 400, `{"origin": [], "destination": ["101", "102"]}`},
		{"invalid-payload-origin-toofew", 400, `{"origin": ["101"], "destination": ["101", "102"]}`},
		{"invalid-payload-origin-toomany", 400, `{"origin": ["101", "103", "104"], "destination": ["101", "102"]}`},
		{"invalid-payload-dest-missing", 400, `{"origin": ["101", "102"]}`},
		{"invalid-payload-dest-empty", 400, `{"destination": [], "origin": ["101", "102"]}`},
		{"invalid-payload-dest-toofew", 400, `{"destination": ["101"], "origin": ["101", "102"]}`},
		{"invalid-payload-dest-toomany", 400, `{"destination": ["101", "103", "104"], "origin": ["101", "102"]}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/orders", strings.NewReader(testCase.payload))
			response := executeRequest(router, req)

			checkResponseCode(t, testCase.code, response.Code)
			checkHeader(t, "application/json; charset=utf-8", response.Header()["Content-Type"])

			switch testCase.code {
			case 200:
				checkBody(t, types.OrderResponse{}, response)
			case 400:
				checkBody(t, types.ErrorResponse{}, response)
			}

			t.Logf("Response body: %s\n", response.Body.String())
		})
	}
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

func checkHeader(t *testing.T, expected string, actual []string) {
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
}

func checkBody(t *testing.T, expected interface{}, request *httptest.ResponseRecorder) {
	jsonErr := json.Unmarshal(request.Body.Bytes(), &expected)
	if jsonErr != nil {
		t.Errorf("Expected body to be valid %T\n", expected)
	}
}
