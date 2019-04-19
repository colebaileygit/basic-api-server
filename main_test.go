package main

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"strings"

	// "io/ioutil"

	// "github.com/go-chi/chi"
	"github.com/gin-gonic/gin"
)

func TestRoutes(t *testing.T) {
	router := Routes()

	testCases := []struct{
		url		string
		command	string
		code	int
		payload string
	}{
		{"/orders", 	"POST", 	200,	`{"origin": ["100", "100"], "destination": ["101", "101"]}`},
		{"/orders/", 	"POST", 	307, 	`{"origin": ["100", "100"], "destination": ["101", "101"]}`},
		{"/orders/1", 	"POST", 	404, 	""},
		// {"/order", 		"POST", 	404, 	""},
		// {"/",			"POST", 	404, 	""},
		// {"/orders", 	"PATCH", 	200, 	""},
		// {"/orders", 	"GET", 		200, 	""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.command + "->" + testCase.url, func(t *testing.T) {
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
