package web

// Importing necessary packages
import (
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
)

// SignUpForm is a function that handles the sign up form submission.
// It parses the form data, creates a new account in the database,
// and updates the cookies.
func SignUpForm(w http.ResponseWriter, r *http.Request) {
	// Initialize an empty account
	var account database.Account

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("pswrd")

	// If the created account is the same as the empty account,
	// execute the home template with the sign up error
	Acc, signUpError, err := database.CreateAccount(email, password, username, false, false)
	if err != nil {
		return
	}

	if Acc == account {
		tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
		tmpl.Execute(w, signUpError)
		return
	}

	// Update the cookies
	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}

	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// LoginForm is a function that handles the login form submission.
// It parses the form data, recovers the account from the database,
// and updates the cookies.
func LoginForm(w http.ResponseWriter, r *http.Request) {
	// Initialize an empty account
	var account database.Account

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	identif := r.Form.Get("identif")
	password := r.Form.Get("pswrd")

	Acc, logInError, err := database.RecoverAccount(identif, password)
	if err != nil {
		// Handle the error
		return
	}

	// If the account found in the database is the same as the empty account,
	// execute the home template with the login error
	if Acc == account {
		tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
		tmpl.Execute(w, logInError)
		return
	}

	// Update the cookies
	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}

	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// LogOutForm is a function that handles the logout form submission.
// It resets the account cookie to a default "Guest" account.
func LogOutForm(w http.ResponseWriter, r *http.Request) {
	// Initialize a default "Guest" account
	var Acc = database.Account{Id: "0", Username: "Guest", ImageUrl: "https://i.pinimg.com/474x/63/bc/94/63bc9469cae29b897565a08f0647db3c.jpg"}

	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}
	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

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

	if username != "QuentinDV" && username != "OwandjiD" {
		if banstatus == "true" {
			db, err := database.ConnectDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.UnBanAccount(db, id)

		} else {
			db, err := database.ConnectDB("database.db")
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

	if username != "QuentinDV" && username != "OwandjiD" {
		if moderator == "true" {
			db, err := database.ConnectDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromModerator(db, id)

		} else {
			db, err := database.ConnectDB("database.db")
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

	if username != "QuentinDV" && username != "OwandjiD" {
		if admin == "true" {
			db, err := database.ConnectDB("database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromAdmin(db, id)
		} else {
			db, err := database.ConnectDB("database.db")
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

	if username != "QuentinDV" && username != "OwandjiD" {
		db, err := database.ConnectDB("database.db")
		if err != nil {
			return
		}
		defer db.Close()
		database.DeleteAccount(db, id)
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func PfpWithUrlForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	imageUrl := r.Form.Get("imageUrl")
	username := r.Form.Get("username")
	fmt.Println(username)

	// Change the image url in the database
	db, err := database.ConnectDB("database.db")
	if err != nil {
		return
	}
	defer db.Close()
	database.ChangeImageUrl(db, id, imageUrl)

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, username)
	if err != nil {
		fmt.Println("Error getting account by username:", err)
		return
	}
	// Update the cookies
	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}

	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/userprofile", http.StatusSeeOther)
}

func PfpWithImageForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(4 << 20) // Set maxMemory parameter to 4MB
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username and id from the form data
	username := r.Form.Get("username")
	id := r.Form.Get("userId")

	// Get the profile picture file from the form data
	file, _, err := r.FormFile("profilePicture")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the profile picture file to the server
	database.SaveFile("./assets/img/pfp/"+id+".png", file)

	db, err := database.ConnectDB("database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Change the image url in the database
	imageUrl := "./assets/img/pfp/" + id + ".png"
	err = database.ChangeImageUrl(db, id, imageUrl)
	if err != nil {
		fmt.Println("Error changing image url:", err)
		return
	}

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, username)
	if err != nil {
		fmt.Println("Error getting account by username:", err)
		return
	}

	// Update the cookies
	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}

	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/userprofile", http.StatusSeeOther)
}
