package orders

import (
	"testing"

	"github.com/colebaileygit/basic-api-server/types"
)

func TestValidatePlaceOrderParams(t *testing.T) {
	testCases := []struct {
		description string
		valid       bool
		origin      []string
		destination []string
	}{
		{"valid-payload", true, []string{"100", "102"}, []string{"103", "104"}},
		{"invalid-payload-origin-toofew", false, []string{"101"}, []string{"101", "102"}},
		{"invalid-payload-origin-toomany", false, []string{"101", "103", "104"}, []string{"101", "102"}},
		{"invalid-payload-origin-empty", false, []string{}, []string{"101", "102"}},
		{"invalid-payload-dest-empty", false, []string{"101", "102"}, []string{}},
		{"invalid-payload-dest-toofew", false, []string{"101", "102"}, []string{"101"}},
		{"invalid-payload-dest-toomany", false, []string{"101", "102"}, []string{"101", "103", "104"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			params := types.PlaceOrderParams{
				Origin:      testCase.origin,
				Destination: testCase.destination,
			}

			if validatePlaceOrderParams(params) != testCase.valid {
				validString := "invalid"
				if testCase.valid {
					validString = "valid"
				}
				t.Errorf("Expected origin: %s and destination: %s to be %s", testCase.origin, testCase.destination, validString)
			}
		})
	}

}
