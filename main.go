package main

import (
	"fmt"
	"forum/assets/go/web"
	"net/http"
)


func main() {
	// Pages
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/login", web.LogIn)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/settings", web.Settings)

	// Elements
	http.Handle("/css", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/settings")
	http.ListenAndServe(":8080", nil)
}
