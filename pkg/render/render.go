package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/almacitunaberk/goforweb/pkg/config"
	"github.com/almacitunaberk/goforweb/pkg/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// Sets the AppConfig for the Render package
func NewTemplates(a *config.AppConfig) {
	app = a;
}

// var templateCache = make(map[string]*template.Template)

// This function adds app-wide shared data to the template
func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	/* OLD WAY of caching before we created AppConfig file
	// Create a Cache for Template
	templateCache, err := createTemplateCache();
	if err != nil {
		log.Fatal(err)
	}
	// Get requested Template from Cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Couldn't find the template")
	}

	// For error handling reasons, instead of writing to the HTTP ResponseWrite object,
	//		we create a byte buffer and try to write it there
	buf := new(bytes.Buffer)

	err = template.Execute(buf, nil)

	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
	*/
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache;
	} else {
		templateCache, _ = CreateTemplateCache();
	}

	template, ok := templateCache[tmpl];

	if !ok {
		log.Fatal("Couldn't find the template")
	}

	buf := new(bytes.Buffer);

	data = AddDefaultData(data, r);

	err := template.Execute(buf, data);

	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w);

	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	// Get all files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return templateCache, err
	}

	// Range through all pages
	for _, page := range pages {
		// First, get the name of the template without the preceding folder names. EX: home.page.tmpl
		name := filepath.Base(page);
		// Then, parse the template file with the given name AND store that parsed template in a template called whatever the value of name is
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}
		// Then, we need to look for the layouts in the app. We want to parse all of them at once so that it is saved into the cache and
		//		can be fetched very quickly
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return templateCache, err
		}
		if len(matches) > 0 {
			// If there are layouts, then parse those layouts AND ADD THEM to templateSet variable
			templateSet, err = templateSet.ParseGlob("./templates/*layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}
		// Set the created templateSet to the cache
		templateCache[name] = templateSet;
	}
	return templateCache, nil
}

/* EASY WAY OF CACHING HTML TEMPLATES, WE WILL USE A MORE COMPLEX BUT BETTER METHOD

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
*/