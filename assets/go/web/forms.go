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

}
