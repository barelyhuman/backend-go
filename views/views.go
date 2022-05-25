package views

import (
	"embed"
	"html/template"
)

//go:embed **/*.html *.html
var embedFS embed.FS

type ViewEngine struct {
	Templates *template.Template
}

// GetTemplates - parse and get all templates
func (ve *ViewEngine) GetTemplates() error {
	allTemplates, err := template.ParseFS(embedFS,
		"**/*.html",
		"*.html",
	)
	if err != nil {
		return err
	}
	ve.Templates = allTemplates
	return nil
}
