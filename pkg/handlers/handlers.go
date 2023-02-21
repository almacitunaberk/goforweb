package handlers

import (
	"net/http"

	"github.com/almacitunaberk/goforweb/pkg/config"
	"github.com/almacitunaberk/goforweb/pkg/models"
	"github.com/almacitunaberk/goforweb/pkg/render"
)

// We are going to use Repository Pattern

// Repo object is used by handlers
var Repo *Repository

// Repository object
type Repository struct {
	App *config.AppConfig
}

// Creating a repository object using AppConfig object
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// QUEST: Why do we need two different funcitons: NewRepo and NewHandlers for setting the Repo from the app?
// ANS:

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}


func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}
