package main

import (
	"fmt"
	"net/http"
)

const PORT = ":8080"


func Home(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "This is home page");
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Bytes written to HTTP Response Writer was %d \n", n));
}


func About(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "This is about page");
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Bytes written to HTTP Response Writer was %d \n", n));
}

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Server listening on port %s", PORT))

	 _ = http.ListenAndServe(PORT, nil);
}