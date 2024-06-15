package web

import (
	"fmt"
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
	fmt.Println(CategoryID)

	db, err := database.ConnectCategoriesDB("db/database.db")
	if err != nil {
		return
	}

	category, err := database.GetCategorybyID(db, CategoryID)
	if err != nil {
		return
	}

	// Execute the user profile template with the Category struct
	tmpl := template.Must(template.ParseFiles("assets/html/category.html"))
	tmpl.Execute(w, category)
}

func LikeForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	PostID := r.Form.Get("LikeID")

	db, err := database.ConnectUserDataDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	IsThisPostLiked := database.IsThisPostLiked(db, ConnectedAccount.Id, PostID)
	if IsThisPostLiked {
		err = database.RemoveLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error removing liked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	} else {
		err = database.RemoveDisLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error removing disliked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}

		err = database.AddLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error adding liked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	}

	// Redirige l'utilisateur vers la page précédente
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	PostID := r.Form.Get("DislikeID")

	db, err := database.ConnectUserDataDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	IsThisPostDisliked := database.IsThisPostDisliked(db, ConnectedAccount.Id, PostID)
	if IsThisPostDisliked {
		err = database.RemoveDisLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error removing disliked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	} else {
		err = database.RemoveLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error removing liked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}

		err = database.AddDisLikedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error adding disliked post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	}

	// Redirige l'utilisateur vers la page précédente
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
