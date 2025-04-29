package main_test

import (
	"encoding/json"
	"os"
	"testing"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

func BenchmarkJsonUnmarshal(b *testing.B) {
	content, err := os.ReadFile("./swagger.json")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		t := spec.Swagger{}
		err := json.Unmarshal(content, &t)
		if err != nil {
			b.Fatal(err)
		}
	}
}
