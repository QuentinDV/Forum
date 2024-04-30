package web

import (
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

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

	// Mettre à jour les cookies
	// Créer un nouveau cookie pour l'Account
	accountCookie := &http.Cookie{
		Name:  "account",
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}
	// Définir le cookie
	http.SetCookie(w, accountCookie)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		// Gérer l'erreur
		return
	}

	if Acc == account {
		tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
		tmpl.Execute(w, logInError)
		return
	}

	fmt.Println("Account Founded")
	// Mettre à jour les cookies
	// Créer un nouveau cookie pour l'Account
	accountCookie := &http.Cookie{
		Name:  "account",
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}

	// Définir le cookie
	http.SetCookie(w, accountCookie)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogOutForm(w http.ResponseWriter, r *http.Request) {
	var Acc = database.Account{Id: "0", Username: "Guest", ImageUrl: "https://i.pinimg.com/474x/63/bc/94/63bc9469cae29b897565a08f0647db3c.jpg"}
	// Créer un nouveau cookie pour l'Account
	accountCookie := &http.Cookie{
		Name:  "account",
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}
	// Définir le cookie
	http.SetCookie(w, accountCookie)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
