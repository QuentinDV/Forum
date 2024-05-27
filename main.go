package main

// Importing necessary packages
import (
	"fmt"
	"forum/assets/go/web"
	"net/http"
)

func main() {
	// Pages
	http.HandleFunc("/", web.LogIn)
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/admin", web.Admin)

	// Forms
	http.HandleFunc("/signupform", web.SignUpForm)
	http.HandleFunc("/loginform", web.LoginForm)
	http.HandleFunc("/guestform", web.LogOutForm)
	http.HandleFunc("/logoutform", web.LogOutForm)
	http.HandleFunc("/banUser", web.BanForm)
	http.HandleFunc("/promoteToModerator", web.ModeratorForm)
	http.HandleFunc("/promoteToAdmin", web.AdminForm)

	// Elements
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("./assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Links
	fmt.Println("\nPlay : http://localhost:8080/admin")
	fmt.Println("\nPlay : http://localhost:8080/categories")
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
