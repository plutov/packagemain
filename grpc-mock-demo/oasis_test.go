package grpcmockdemo

import (
	"fmt"
	"testing"

	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
	"google.golang.org/grpc/examples/route_guide/routeguide"
)

func TestIsOasis(t *testing.T) {
	routeGuideSvcAddr := "localhost:50001"
	srv := createMockGrpcServer(routeGuideSvcAddr)
	defer srv.Close()

	tests := []struct {
		point Point
		want  bool
	}{
		{
			Point{10, 10},
			true,
		},
		{
			Point{20, 20},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Point(%d,%d)", tt.point.latitude, tt.point.longitude), func(t *testing.T) {
			got, err := IsOasis(t.Context(), routeGuideSvcAddr, tt.point)
			if err != nil {
				t.Fatalf("not expected err: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func createMockGrpcServer(addr string) *grpcmock.Server {
	srv := grpcmock.NewServer(
		grpcmock.WithAddress(addr),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.RegisterService(routeguide.RegisterRouteGuideServer),
		func(s *grpcmock.Server) {
			s.ExpectUnary(routeguide.RouteGuide_GetFeature_FullMethodName).
				WithPayload(&routeguide.Point{
					Latitude:  10,
					Longitude: 10,
				}).
				Return(&routeguide.Feature{
					Name: FeatureForest,
				})

			s.ExpectUnary(routeguide.RouteGuide_GetFeature_FullMethodName).
				UnlimitedTimes().
				Return(&routeguide.Feature{
					Name: FeatureDesert,
				})
		},
	)

	return srv
}
