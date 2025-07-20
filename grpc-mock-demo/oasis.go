package grpcmockdemo

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
	// give some deadline to context
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// create a conn to grpc server
	conn, err := grpc.NewClient(routeGuideSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err
	}

	client := routeguide.NewRouteGuideClient(conn)

	feature, err := client.GetFeature(ctx, &routeguide.Point{
		Latitude:  point.latitude,
		Longitude: point.longitude,
	})
	if err != nil || feature.Name != FeatureForest {
		return false, err
	}

	// 0 0 0
	// 0 1 1
	// 0 0 0

	dx := []int32{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int32{-1, -1, -1, 0, 0, 1, 1, 1}
	for i := range dx {
		feature, err := client.GetFeature(ctx, &routeguide.Point{
			Latitude:  point.latitude + dx[i],
			Longitude: point.longitude + dy[i],
		})
		if err != nil || feature.Name != FeatureDesert {
			return false, err
		}
	}

	return true, nil
}
