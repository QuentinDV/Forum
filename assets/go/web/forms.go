package web

import (
	"net/http"
)

func SignUpForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	// username := r.Form.Get("username")
	load := r.Form.Get("load")
	delete := r.Form.Get("delete")

	if load == "true" {
		http.Redirect(w, r, "/hangman/load", http.StatusSeeOther)
	} else if delete == "true" {
		http.Redirect(w, r, "/hangman/delete", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/hangman/choose_difficulty", http.StatusSeeOther)
}
