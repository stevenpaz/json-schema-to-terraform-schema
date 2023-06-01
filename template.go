package main

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/tf-schema.go.tmpl
var schemaTemplate embed.FS

type Struct struct {
	Name   string
	Fields map[string]string
}

type TemplateData struct {
	FileName string
	Structs  []Struct
}

func GetTemplate() (*template.Template, error) {
	// Process the template
	tmpl, err := template.ParseFS(schemaTemplate, "templates/tf-schema.go.tmpl")
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	return tmpl, nil
}
