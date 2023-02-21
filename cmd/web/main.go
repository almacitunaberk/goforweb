package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/almacitunaberk/goforweb/pkg/config"
	"github.com/almacitunaberk/goforweb/pkg/handlers"
	"github.com/almacitunaberk/goforweb/pkg/render"
)

const PORT = ":8080"



func main() {
	var app config.AppConfig
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache");
	}
	app.TemplateCache = templateCache
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	render.NewTemplates(&app)
	handlers.NewHandlers(repo)
	/*
	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	*/
	fmt.Println(fmt.Sprintf("Server listening on port %s", PORT))

	// _ = http.ListenAndServe(PORT, nil)
	srv := &http.Server{
		Addr: PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}