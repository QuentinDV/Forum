package main

import "forum/assets/go/database"

func main() {
	database.CreateAccount("quentin.dassivignon@ynov.com", "quentin123", "blazefast")

	// // Pages
	// http.HandleFunc("/home", web.Home)
	// http.HandleFunc("/categories", web.Categories)
	// http.HandleFunc("/login", web.LogIn)
	// http.HandleFunc("/signup", web.SignUp)
	// http.HandleFunc("/settings", web.Settings)

	// // Elements
	// http.Handle("/css", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	// http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// // Liens
	// fmt.Println("\nPlay : http://localhost:8080/home")
	// fmt.Println("\nPlay : http://localhost:8080/login")
	// http.ListenAndServe(":8080", nil)
}
