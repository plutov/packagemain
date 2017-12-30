//go:generate go-bindata -o tpl.go tpl
package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	b, err := Asset("tpl/page.html")
	if err != nil {
		log.Fatalf("unable to get template: %v", err)
	}
	t, err := template.New("base").Parse(string(b))
	if err != nil {
		log.Fatalf("unable to parse template: %v", err)
	}
	err = t.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("unable to execute template: %v", err)
	}
}
