package orders

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"googlemaps.github.io/maps"

	"github.com/colebaileygit/basic-api-server/types"
)

func CalculateDistance(params types.PlaceOrderParams) (distance int, err error) {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Panicf("Google Maps client could not be initialized: %s\n", err)
	}

	request := &maps.DistanceMatrixRequest{
		Origins:      []string{strings.Join(params.Origin, ",")},
		Destinations: []string{strings.Join(params.Destination, ",")},
	}

	// log.Printf("Request: %s | %+v\n", request.Origins[0], request)

	distanceMatrix, err := client.DistanceMatrix(context.Background(), request)

	if err != nil {
		return 0, fmt.Errorf("Google Maps client could not return for query: %+v \n%s\n", request, err)
	}

	if distanceMatrix.Rows == nil || len(distanceMatrix.Rows) == 0 || len(distanceMatrix.Rows[0].Elements) == 0 {
		return 0, fmt.Errorf("Google Maps result did not contain distance value: %+v\n", distanceMatrix)
	}

	// log.Printf("Distances: %+v\n", distanceMatrix.Rows[0].Elements[0])

	element := distanceMatrix.Rows[0].Elements[0]

	if element.Status != "OK" {
		return 0, fmt.Errorf("Google Maps could not return distance for query: %s", element.Status)
	}

	result := element.Distance.Meters

	if result <= 0 {
		return result, errors.New("Route must have valid distance (>0) between origin and destination.")
	}

	return result, nil
}
