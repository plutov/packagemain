//go:generate go-bindata -o tpl.go tpl
package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

func main() {
	var useGB bool
	flag.BoolVar(&useGB, "go-bindata", false, "")
	flag.Parse()

	var t *template.Template
	if useGB {
		b, err := Asset("tpl/page.html")
		if err != nil {
			log.Fatalf("unable to asset template: %v", err)
		}
		t, err = template.New("").Parse(string(b))
		if err != nil {
			log.Fatalf("unable to parse template: %v", err)
		}
	} else {
		t = template.Must(template.ParseFiles("tpl/page.html"))
	}

	if t != nil {
		err := t.ExecuteTemplate(os.Stdout, "base", nil)
		if err != nil {
			log.Fatalf("unable to execute template: %v", err)
		}
	}
}
