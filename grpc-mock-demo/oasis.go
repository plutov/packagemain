package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/route_guide/routeguide"
)

type Point struct {
	latitude  int32
	longitude int32
}

const (
	FeatureForest = "forest"
	FeatureDesert = "desert"
)

func IsOasis(ctx context.Context, routeGuideSvcAddr string, point Point) (bool, error) {
	// give some deadline to the context
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	// create a connection to the server
	conn, err := grpc.NewClient(routeGuideSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}

	// create a client using the connection
	client := routeguide.NewRouteGuideClient(conn)

	// check if the current point is a forest
	feature, err := client.GetFeature(ctx, &routeguide.Point{
		Latitude:  point.latitude,
		Longitude: point.longitude,
	})
	if err != nil {
		return false, err
	}
	if feature.Name != FeatureForest {
		return false, nil
	}

	// check if all adjacent points are deserts (with diagonal checks)
	// these two slices represent the relative coordinates of the adjacent points
	dLat := []int32{-1, 0, 1, -1, 1, -1, 0, 1}
	dLong := []int32{-1, -1, -1, 0, 0, 1, 1, 1}
	for i := 0; i < len(dLat); i++ {
		adjacentPoint := &routeguide.Point{
			Latitude:  point.latitude + dLat[i],
			Longitude: point.longitude + dLong[i],
		}
		feature, err := client.GetFeature(ctx, adjacentPoint)
		if err != nil {
			return false, err
		}
		if feature.Name != FeatureDesert {
			return false, nil
		}
	}

	// if all checks passed, return true
	return true, nil
}
