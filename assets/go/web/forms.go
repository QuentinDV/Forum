package web

import (
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

type SignUpError struct {
}

func SignUpForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("pswrd")

	// Acc := database.Account{Username: username, Email: email, Password: password}
	Acc, err := database.CreateAccount(email, password, username)
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles("assets/html/.html"))
	tmpl.Execute(w, Acc)
}
