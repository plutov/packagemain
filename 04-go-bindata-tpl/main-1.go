package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	t := template.Must(template.ParseFiles("tpl/page.html"))

	err := t.ExecuteTemplate(os.Stdout, "base", nil)
	if err != nil {
		log.Fatalf("unable to execute template: %v", err)
	}
}
