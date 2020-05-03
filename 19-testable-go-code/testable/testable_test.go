package testable_test

import (
	"testing"

	"github.com/plutov/packagemain/drafts/19-testable-go-code/testable"
)

func TestStrInSlice(t *testing.T) {
	var tests = []struct {
		slice []string
		find  string
		want  bool
	}{
		{[]string{"a", "b"}, "c", false},
		{[]string{"a", "b"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.find, func(t *testing.T) {
			got := testable.StrInSlice(tt.slice, tt.find)
			if got != tt.want {
				t.Errorf("expecting %t, got %t", tt.want, got)
			}
		})
	}
}

func TestGetAverageStarsPerRepo(t *testing.T) {
	var tests = []struct {
		username string
		want     float64
	}{
		{"octocat", 4},
	}

	mock := new(testable.Mock)
	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			got, err := testable.GetAverageStarsPerRepo(mock, tt.username)
			if err != nil {
				t.Errorf("expecting nil err, got %v", err)
			}
			if got != tt.want {
				t.Errorf("expecting %f, got %f", tt.want, got)
			}
		})
	}
}
