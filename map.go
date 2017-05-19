package main

import (
	"errors"
	"os"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

const (
	MaxRetry = 3
)

// calculate the duration and distance from Google and update the record with latest value
func CalcRoute(token string, record Record) (err error) {
	var routes *maps.DistanceMatrixResponse

	for i := 0; i < MaxRetry; i++ {
		routes, err = GetShortestByGoogle(record.Path)

		// break the retry mechanism if there is no error
		if err == nil {
			break
		}

		// sleep for a while to avoid hitting Google too hard
		time.Sleep(time.Millisecond * 500)
	}

	// calculate the total duration and distance
	for _, route := range routes.Rows {
		for _, ele := range route.Elements {
			record.TotalTime += int(ele.Duration.Seconds())
			record.TotalDistance += ele.Distance.Meters
		}
	}

	// mark the record is success
	record.Status = StatusSuccess

	// update the database with the latest value from Google
	err = db.Set(token, record, 0).Err()
	if err != nil {
		// we can send an alert to log this error
		// if the data cannot insert into database
		return
	}

	return nil
}

// send a HTTP request to google and get the shortest path (DistanceMatrix)
func GetShortestByGoogle(points Path) (routes *maps.DistanceMatrixResponse, err error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		return nil, errors.New("'GOOGLE_API_KEY' is missing")
	}

	// Get the api key from environmental variable instead of hard code it
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return
	}

	// shift the first element from the array
	origin, points := points[0], points[1:]
	destination := []string{}

	// create an array to store formatted lat log
	for _, dest := range points {
		destination = append(destination, dest[0]+","+dest[1])
	}

	r := &maps.DistanceMatrixRequest{
		Origins: []string{
			origin[0] + "," + origin[1],
		},
		Destinations: destination,
	}

	resp, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
