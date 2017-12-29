package main

import (
	"html/template"
	"os"
	"testing"
)

func BenchFile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := template.Must(template.ParseFiles("tpl/page.html"))
		t.ExecuteTemplate(os.Stdout, "base", nil)
	}
}

func BenchBindata(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b, _ := Asset("tpl/page.html")
		t, _ := template.New("").Parse(string(b))
		t.ExecuteTemplate(os.Stdout, "base", nil)
	}
}
