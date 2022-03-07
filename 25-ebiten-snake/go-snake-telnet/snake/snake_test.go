// Copyright (c) 2017 Alex Pliutau
package snake

import (
	"testing"
)

func BenchmarkMove(b *testing.B) {
	s := initialSnake()

	for n := 0; n < b.N; n++ {
		err := s.move()
		if err != nil {
			b.Fatalf("failed to move: %v", err)
		}
	}
}

func BenchmarkHits(b *testing.B) {
	s := initialSnake()

	for n := 0; n < b.N; n++ {
		s.hits(coord{
			x: 1,
			y: 1,
		})
	}
}
