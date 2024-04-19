package main

import (
	"fmt"
	"net/http"
)

// Page Principale du Forum
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/home.html")
}

func Categories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/categories.html")
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/signup.html")
}

func Settings(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/settings.html")
}

func main() {
	// Pages
	http.HandleFunc("/home", Home)
	http.HandleFunc("/categories", Home)
	http.HandleFunc("/login", Home)
	http.HandleFunc("/signup", Home)
	http.HandleFunc("/settings", Home)

	// Elements
	http.Handle("/css", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/settings")
	http.ListenAndServe(":8080", nil)
}
