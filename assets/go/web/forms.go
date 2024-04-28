package web

import (
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

// var ConnectedAccount = database.Account{Id: "0", Username: "Guest", ImageUrl: "https://i.pinimg.com/474x/63/bc/94/63bc9469cae29b897565a08f0647db3c.jpg"}

func SignUpForm(w http.ResponseWriter, r *http.Request) {
	var account database.Account
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("pswrd")

	fmt.Println(email, password, username)

	Acc, signUpError, err := database.CreateAccount(email, password, username, false)
	if err != nil {
		return
	}
	
	if Acc == account {
		tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
		tmpl.Execute(w, signUpError)
		return
	}

	tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
	tmpl.Execute(w, Acc)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	var account database.Account
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	identif := r.Form.Get("identif")
	password := r.Form.Get("pswrd")

	fmt.Println(identif, password)

	Acc, logInError, err := database.RecoverAccount(identif, password)
	if err != nil {
		return
	}

	if Acc == account {
		tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
		tmpl.Execute(w, logInError)
		return
	}

	tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
	tmpl.Execute(w, Acc)
}
