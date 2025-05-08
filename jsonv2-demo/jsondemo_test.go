package main

import (
	json "encoding/json/v2"
	"os"
	"testing"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

func BenchmarkJsonUnmarshal(b *testing.B) {
	f, err := os.ReadFile("./swagger.json")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		out := spec.Swagger{}
		if err := json.Unmarshal(f, &out); err != nil {
			b.Fatal(err)
		}
	}
}
