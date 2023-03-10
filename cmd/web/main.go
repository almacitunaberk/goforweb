package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/almacitunaberk/goforweb/pkg/config"
	"github.com/almacitunaberk/goforweb/pkg/handlers"
	"github.com/almacitunaberk/goforweb/pkg/helpers"
	"github.com/almacitunaberk/goforweb/pkg/models"
	"github.com/almacitunaberk/goforweb/pkg/render"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	// What we will be putting into the session:
	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO|\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR|\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Configuring sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// In production, the next should be true
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Cannot create template cache");
		log.Println(err)
		return
	}
	app.TemplateCache = templateCache
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	render.NewTemplates(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
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