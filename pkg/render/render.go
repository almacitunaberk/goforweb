package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templateCache = make(map[string]*template.Template)

func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.tmpl");
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		return;
	}
}

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template;
	var err error

	// If we already have the template in cache, we want to use that template
	_, inMap := templateCache[t];
	if !inMap {
		// We need to create the template
		log.Println("Creating the template")
		err = createTemplateCache(t);
		if err != nil {
			log.Println(err)
		}
	} else {
		// We already have the template in the memory
		log.Println("Cache hit for the template")
	}
	// At this point, we either created a template or it was already created. Thus, templateCache[t] definitely exists
	tmpl = templateCache[t];
	err = tmpl.Execute(w, nil);
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {
	templates := []string {
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}
	// Parse the template
	tmpl, err := template.ParseFiles(templates...);
	if err != nil {
		return err
	}
	// Add template to cache
	templateCache[t] = tmpl;
	return nil;
}