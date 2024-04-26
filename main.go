package main

import (
	"fmt"
	"forum/assets/go/database"
	"forum/assets/go/web"
	"net/http"
)

var ConnectedAccount = database.Account{Id: "0", Username: "Guest", ImageUrl: "https://i.pinimg.com/474x/63/bc/94/63bc9469cae29b897565a08f0647db3c.jpg"}

func main() {
	// Pages
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/login", web.LogIn)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/settings", web.Settings)

	//Forms
	http.HandleFunc("/signupform", web.SignUpForm)

	// Elements
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("./assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/settings")
	fmt.Println("\nPlay : http://localhost:8080/categories")
	fmt.Println("\nPlay : http://localhost:8080/login")
	fmt.Println("\nPlay : http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}
