package web

import (
	"net/http"
)

// Page Principale du Forum
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/html/home.html")
}

func Categories(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/html/categories.html")
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/html/login.html")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/html/signup.html")
}

func Settings(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/html/settings.html")
}
