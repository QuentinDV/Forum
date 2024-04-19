package main

import (
	"fmt"
	"net/http"
)

// Page Principale du Forum
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/home.html")
}

func main() {
	// Pages
	http.HandleFunc("/home", Home)

	// Elements
	http.Handle("/assets/static/", http.StripPrefix("/assets/static/", http.FileServer(http.Dir("./assets/static"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}
