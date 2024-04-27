package web

import (
	"fmt"
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

	fmt.Println(email, password, username)

	Acc, err := database.CreateAccount(email, password, username, false)
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
	tmpl.Execute(w, Acc)
}
