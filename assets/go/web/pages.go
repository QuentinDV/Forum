package web

// Importing necessary packages
import (
	"database/sql"
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

// Home is the main page of the forum.
func Home(w http.ResponseWriter, r *http.Request) {
	// Serve the home page
	http.ServeFile(w, r, "assets/html/home.html")
}

// Categories page of the forum.
func Categories(w http.ResponseWriter, r *http.Request) {
	// Serve the categories page
	http.ServeFile(w, r, "assets/html/categories.html")
}

// Categories page of the forum.
func UserProfile(w http.ResponseWriter, r *http.Request) {
	// Serve the categories page
	http.ServeFile(w, r, "assets/html/userprofile.html")
}

// LogIn page of the forum.
func LogIn(w http.ResponseWriter, r *http.Request) {
	// Serve the login page
	http.ServeFile(w, r, "assets/html/login.html")
}

// SignUp page of the forum.
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Serve the signup page
	http.ServeFile(w, r, "assets/html/signup.html")
}

// Admin page of the forum.
func Admin(w http.ResponseWriter, r *http.Request) {
	// Open the database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	allAcc, err := database.GetAllAccounts(db)
	if err != nil {
		fmt.Println("Error getting all accounts:", err)
		return
	}

	// Serve the admin page
	tmpl := template.Must(template.ParseFiles("assets/html/admin.html"))
	tmpl.Execute(w, allAcc)
}
