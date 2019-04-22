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
		{"valid-payload", true, []string{"22", "102"}, []string{"23", "104"}},
		{"valid-payload", true, []string{"-22", "102"}, []string{"23", "104"}},
		{"invalid-payload-format", false, []string{"22,3", "52"}, []string{"24", "53"}},
		{"invalid-payload-origin-toofew", false, []string{"22"}, []string{"22", "102"}},
		{"invalid-payload-origin-toomany", false, []string{"22", "103", "104"}, []string{"22", "102"}},
		{"invalid-payload-origin-empty", false, []string{}, []string{"22", "102"}},
		{"invalid-payload-dest-empty", false, []string{"22", "102"}, []string{}},
		{"invalid-payload-dest-toofew", false, []string{"22", "102"}, []string{"22"}},
		{"invalid-payload-dest-toomany", false, []string{"22", "102"}, []string{"22", "103", "104"}},
		{"invalid-payload-lat-bounds", false, []string{"-91", "100"}, []string{"-90", "99"}},
		{"invalid-payload-lat-bounds", false, []string{"90", "100"}, []string{"91", "99"}},
		{"invalid-payload-lng-bounds", false, []string{"53", "180"}, []string{"53", "181"}},
		{"invalid-payload-lng-bounds", false, []string{"53", "-181"}, []string{"53", "-180"}},
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
