package web

import (
	"forum/assets/go/database"
	"net/http"
	"strings"
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

func getAccountfromCookie(r *http.Request) database.Account {
	ConnectedAccount := database.Account{}
	cookie, err := r.Cookie("account")
	if err != nil {
		return ConnectedAccount
	}
	cookieValue := cookie.Value
	cookieValues := strings.Split(cookieValue, "|")
	return database.Account{
		Id:           cookieValues[0],
		Email:        cookieValues[1],
		Password:     cookieValues[2],
		Username:     cookieValues[3],
		ImageUrl:     cookieValues[4],
		IsBan:        cookieValues[5] == "true",
		IsAdmin:      cookieValues[6] == "true",
		CreationDate: cookieValues[7],
	}
}
