package testable

import "testing"

func Test_strInSlice(t *testing.T) {
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
			got := strInSlice(tt.slice, tt.find)
			if got != tt.want {
				t.Errorf("expecting %t, got %t", tt.want, got)
			}
		})
	}
}

func Test_getAveragePerRepo(t *testing.T) {
	var tests = []struct {
		username string
		want     float64
	}{
		{"octocat", 1480.375000},
		{"plutov", 15.566667},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			got, err := getAveragePerRepo(tt.username)
			if err != nil {
				t.Errorf("expecting nil err, got %v", err)
			}
			if got != tt.want {
				t.Errorf("expecting %f, got %f", tt.want, got)
			}
		})
	}
}
