package templates

import (
	"embed"
	"log"
	"text/template"
)

//go:embed *
var FS embed.FS

func ParseTemplate() *template.Template {
    t, err := template.ParseFS(FS, "home.gohtml")
    if err != nil {
        log.Fatal(err)
    }
    return t
}