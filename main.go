package main

import (
	"fmt"
	"forum/assets/go/database"
	"forum/assets/go/web"
	"net/http"
)

func main() {
	database.CreateAccount("quentin.dassivignon@ynov.com", "quentin123", "quentindv", true)
	database.CreateAccount("owandji.dieng@ynov.com", "owandji123", "owandji", true)

	// Pages
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/login", web.LogIn)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/settings", web.Settings)

	// Elements
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("./assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/categories")
	fmt.Println("\nPlay : http://localhost:8080/login")
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/login")
	http.ListenAndServe(":8080", nil)
}
