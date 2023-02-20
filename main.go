package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const PORT = ":8080"


func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.tmpl")
}


func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.tmpl")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl);
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		return;
	}

}

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Server listening on port %s", PORT))

	 _ = http.ListenAndServe(PORT, nil);
}