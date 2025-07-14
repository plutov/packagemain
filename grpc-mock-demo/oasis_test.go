package main

import (
	"fmt"
	"testing"

	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
	"google.golang.org/grpc/examples/route_guide/routeguide"
)

func TestIsOasis(t *testing.T) {
	routeGuideSvcAddr := "localhost:50051"
	srv := createMockRouteGuideServer(routeGuideSvcAddr)
	t.Cleanup(func() {
		srv.Close()
	})

	tests := []struct {
		point Point
		want  bool
	}{
		{
			Point{latitude: 10, longitude: 10},
			true,
		},
		{
			Point{latitude: 20, longitude: 20},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Point(%d,%d)", tt.point.latitude, tt.point.longitude), func(t *testing.T) {
			got, err := IsOasis(t.Context(), routeGuideSvcAddr, tt.point)
			if err != nil {
				t.Fatalf("IsOasis() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("IsOasis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMockRouteGuideServer(addr string) *grpcmock.Server {
	srv := grpcmock.NewServer(
		grpcmock.WithAddress(addr),
		grpcmock.RegisterService(routeguide.RegisterRouteGuideServer),
		grpcmock.WithPlanner(planner.FirstMatch()),
		func(s *grpcmock.Server) {
			// 10,10 is a forest
			s.ExpectUnary(routeguide.RouteGuide_GetFeature_FullMethodName).
				UnlimitedTimes().
				WithPayload(&routeguide.Point{
					Latitude:  10,
					Longitude: 10,
				}).
				Return(&routeguide.Feature{
					Name: "forest",
				})

			// everything else is a desert
			s.ExpectUnary(routeguide.RouteGuide_GetFeature_FullMethodName).
				UnlimitedTimes().
				Return(&routeguide.Feature{
					Name: "desert",
				})
		},
	)

	return srv
}
