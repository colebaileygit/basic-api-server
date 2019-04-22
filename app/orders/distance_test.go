package orders

import (
	"fmt"
	"math"
	"testing"

	"github.com/colebaileygit/basic-api-server/types"
)

func TestCalculateDistance(t *testing.T) {
	testCases := []struct {
		originLat         string
		originLng         string
		destinationLat    string
		destinationLng    string
		approximateMeters int
		hasError          bool
	}{
		{"22.2789632", "114.1858304", "22.316552881533134", "114.21683948898317", 12600, false},
		{"22.316552881533134", "114.21683948898317", "22.315088364385748", "113.9370791205622", 39990, false},
		{"101", "101", "102", "102", 0, true},
		{"22.316552881533134", "114.21683948898317", "22.316552881533134", "114.21683948898317", 0, true}, // Same start + end
		{"22.316552881533134", "114.21683948898317", "49.23606781069188", "117.33454026596814", 0, true},  // Distance too far
		{"-90", "-360", "-91", "-360", 0, true},
	}

	for _, testCase := range testCases {
		description := fmt.Sprintf("(%s,%s),(%s,%s)", testCase.originLat, testCase.originLng, testCase.destinationLat, testCase.destinationLng)
		t.Run(description, func(t *testing.T) {
			params := types.PlaceOrderParams{
				Origin:      []string{testCase.originLat, testCase.originLng},
				Destination: []string{testCase.destinationLat, testCase.destinationLng},
			}
			meters, err := CalculateDistance(params)

			if testCase.hasError {
				if err == nil {
					t.Errorf("Distance queried for request did not error. Distance calculated: %d", meters)
				}
			} else {
				if err != nil {
					t.Errorf("Error while retrieving distance: %s", err)
				}
				if math.Abs(float64(meters-testCase.approximateMeters)) > 50 {
					t.Errorf("Distance was not close to estimated value. Expected: %d, Actual: %d", testCase.approximateMeters, meters)
				}
			}

		})
	}
}
