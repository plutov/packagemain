package teereeder

import (
	"strings"
	"testing"
)

func BenchmarkValidateAndSaveCSVFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		csvFile := strings.NewReader("col1,col2\ncell1,cell2\ncell3,cell4")
		if err := ValidateAndSaveCSVFile(csvFile); err != nil {
			b.Fatal("unable to validate and save csv file")
		}
	}
}

func BenchmarkValidateAndSaveCSVFile2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		csvFile := strings.NewReader("col1,col2\ncell1,cell2\ncell3,cell4")
		if err := ValidateAndSaveCSVFile2(csvFile); err != nil {
			b.Fatal("unable to validate and save csv file")
		}
	}
}
