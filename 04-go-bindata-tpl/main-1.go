package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	t, err := template.ParseFiles("tpl/page.html")
	if err != nil {
		log.Fatalf("unable to parse template: %v", err)
	}

	err = t.ExecuteTemplate(os.Stdout, "base", nil)
	if err != nil {
		log.Fatalf("unable to execute template: %v", err)
	}
}
