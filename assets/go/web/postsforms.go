package web

import (
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

func CategoryForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	CategoryID := r.Form.Get("CategoryID")

	db, err := database.ConnectCategoriesDB("database.db")
	if err != nil {
		return
	}

	category, err := database.GetCategory(db, CategoryID)
	if err != nil {
		return
	}

	// Execute the user profile template with the Category struct
	tmpl := template.Must(template.ParseFiles("assets/html/category.html"))
	tmpl.Execute(w, category)
}
