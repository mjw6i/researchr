package main

import (
	"html/template"
	"os"
)

func main() {
	layout := template.Must(template.ParseFiles("index.htm"))

	layout.Execute(os.Stdout, "")
}
