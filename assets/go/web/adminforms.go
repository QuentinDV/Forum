package web

// Importing necessary packages
import (
	"forum/assets/go/database"
	"log"
	"net/http"
)

func BanForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")
	banstatus := r.Form.Get("banstatus")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if banstatus == "true" {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.UnBanAccount(db, id)

		} else {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.BanAccount(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ModeratorForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")
	moderator := r.Form.Get("moderator")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if moderator == "true" {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromModerator(db, id)

		} else {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.PromoteToModerator(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	admin := r.Form.Get("admin")
	username := r.Form.Get("username")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if admin == "true" {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromAdmin(db, id)
		} else {
			db, err := database.ConnectUserDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.PromoteToAdmin(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteAccountForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		db, err := database.ConnectUserDB("database.db")
		if err != nil {
			return
		}
		defer db.Close()
		database.DeleteAccount(db, id)
	}

	userdb, err := database.ConnectUserDataDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer userdb.Close()
	err = database.DeleteUserData(userdb, id)
	if err != nil {
		log.Fatal(err)
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
