package handlers

import (
	"embed"
	"html/template"
	"log"
)

var templates *template.Template

// InitHandlers 初始化handlers包
func InitHandlers(content embed.FS) {
	var err error
	templates, err = template.ParseFS(content, "static/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
}
