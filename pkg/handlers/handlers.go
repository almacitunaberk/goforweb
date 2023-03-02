package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/almacitunaberk/goforweb/pkg/config"
	"github.com/almacitunaberk/goforweb/pkg/driver"
	"github.com/almacitunaberk/goforweb/pkg/forms"
	"github.com/almacitunaberk/goforweb/pkg/helpers"
	"github.com/almacitunaberk/goforweb/pkg/models"
	"github.com/almacitunaberk/goforweb/pkg/render"
	"github.com/almacitunaberk/goforweb/pkg/repository"
	"github.com/almacitunaberk/goforweb/pkg/repository/dbrepo"
	"github.com/go-chi/chi"
)

// We are going to use Repository Pattern

// Repo object is used by handlers
var Repo *Repository

// Repository object
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// Creating a repository object using AppConfig object
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// QUEST: Why do we need two different funcitons: NewRepo and NewHandlers for setting the Repo from the app?
// ANS:

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}


func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl",&models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl",  &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from the session"))
		return
	}

	sd := res.StartDate.Format("2006-01-01")
	ed := res.EndDate.Format("2006-01-01")

	room, err := m.DB.GetRoomByID(res.RoomID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	// Because we need the following information on the Reservation Summary page
	// 		we will store it in the session
	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	reservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room_restriction := models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: reservationId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(room_restriction)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// Converting Start and End from String to time.Time data type
	layout := `2006-01-02`
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK bool `json:"ok"`
	Message string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	RoomID string `json:"room_id"`
}
func (m *Repository) AvailabilityJSON(w http.ResponseWriter ,r *http.Request) {
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2020-01-01"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK: available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can not get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can not get reservation from session")
		http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
		return
	}
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}


func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}


func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	// Reservation object (res) comes from the session but when the user searches for availability, there
	//		is no specific room selected yet. Thus, we need to assign the room ID to the res object
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// Getting id, s and e params from the URL:
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"));
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	layout := "2020-02-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w,r,"login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w,r,"login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	var email string
	var password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w,r,"/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w,r,"/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w,r,"/", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	return;
}