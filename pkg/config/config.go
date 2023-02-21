package config

import "html/template"

// AppConfig holds app-wide shared data
type AppConfig struct {
	UseCache bool
	TemplateCache map[string]*template.Template
}