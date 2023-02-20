package main

import (
	"fmt"
	"net/http"

	"github.com/almacitunaberk/goforweb/pkg/handlers"
)

const PORT = ":8080"



func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Server listening on port %s", PORT))

	 _ = http.ListenAndServe(PORT, nil);
}