package web

import (
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
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
	// var account database.Account

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}
	identif := r.Form.Get("identif")
	password := r.Form.Get("pswrd")

	Acc, LogInError, err := database.RecoverAccount(identif, password)
	fmt.Println("LogInError:", LogInError)
	if err != nil {
		return
	}
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// If the account found in the database is the same as the empty account,
	// execute the home template with the login error
	// if Acc == account {
	// 	tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
	// 	tmpl.Execute(w, logInError)
	// 	return
	// }

	// fmt.Println("account:", account)
	// fmt.Println("Acc:", Acc)
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
	http.Redirect(w, r, "/home", http.StatusFound)
}

// LogOutForm is a function that handles the logout form submission.
// It resets the account cookie to a default "Guest" account.
func LogOutForm(w http.ResponseWriter, r *http.Request) {
	// Initialize a default "Guest" account
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()
	Acc, err := database.GetAccountbyID(db, "0")
	if err != nil {
		fmt.Println("Error getting account by ID:", err)
		return
	}

	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}
	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GuestForm is a function that handles the guest form submission.
// It sets the account cookie to a default "Guest" account.
func GuestForm(w http.ResponseWriter, r *http.Request) {
	// Initialize a default "Guest" account
	db, err := database.ConnectUserDB("db/database.db")

	if err != nil {
		return
	}
	defer db.Close()
	Acc, err := database.GetAccountbyID(db, "0")
	if err != nil {
		fmt.Println("Error getting account by ID:", err)
		return
	}

	// Create a new cookie for the account
	accountCookie := &http.Cookie{
		Name: "account",
		// The value of the cookie is a string that contains the account's information separated by "|"
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
	}
	// Set the cookie
	http.SetCookie(w, accountCookie)

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func PfpWithImageForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(20 << 20) // Set maxMemory parameter to 20MB
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Get the profile picture file from the form data
	file, _, err := r.FormFile("profilePicture")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the profile picture file to the server
	database.SaveFile("./assets/img/pfp/"+ConnectedAccount.Id+".png", file)

	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Change the image url in the database
	imageUrl := "./assets/img/pfp/" + ConnectedAccount.Id + ".png"
	err = database.ChangeImageUrl(db, ConnectedAccount.Id, imageUrl)
	if err != nil {
		fmt.Println("Error changing image url:", err)
		return
	}

	// Get the account from the database
	Acc, err := database.GetAccountbyID(db, ConnectedAccount.Id)
	if err != nil {
		fmt.Println("Error getting account by ID:", err)
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

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func ChangePwForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve the connected account from the cookie or session
	ConnectedAccount := RetrieveAccountfromCookie(r)

	newpassword := r.Form.Get("newPw")

	// Change the password in the database
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = database.ChangePassword(db, ConnectedAccount.Id, newpassword)
	if err != nil {
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	// Get the updated account from the database
	Acc, err := database.GetAccountByUsername(db, ConnectedAccount.Username)
	if err != nil {
		fmt.Println("Error getting account by username:", err)
		return
	}

	// Update the cookies
	accountCookie := &http.Cookie{
		Name:  "account",
		Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", Acc.Id, Acc.Email, Acc.Password, Acc.Username, Acc.ImageUrl, Acc.IsBan, Acc.IsModerator, Acc.IsAdmin, Acc.CreationDate),
		Path:  "/",
	}
	http.SetCookie(w, accountCookie)

	// Redirect to the user profile page
	http.Redirect(w, r, "/userprofile", http.StatusSeeOther)
}

func RetrieveAccountfromCookie(r *http.Request) database.Account {
	// Retrieve the account cookie
	accountCookie, err := r.Cookie("account")
	if err != nil {
		// If there is an error retrieving the cookie, handle it accordingly
		// For example, redirect to the login page or display an error message
		return database.Account{}
	}

	// Split the cookie value into separate fields
	fields := strings.Split(accountCookie.Value, "|")

	// Create a new account using the cookie data
	account := database.Account{
		Id:           fields[0],
		Email:        fields[1],
		Password:     fields[2],
		Username:     fields[3],
		ImageUrl:     fields[4],
		IsBan:        fields[5] == "true",
		IsModerator:  fields[6] == "true",
		IsAdmin:      fields[7] == "true",
		CreationDate: fields[8],
	}

	return account
}

func AddViewForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the post ID from the form data
	postID := r.Form.Get("PostID")

	db, err := database.ConnectPostDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Increment the number of views of the post
	err = database.IncrementViewtoDB(db, postID)
	if err != nil {
		fmt.Println("Error incrementing views:", err)
		return
	}

	// Redirect to the post page
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// Reset Pfp to default
func ResetPfpForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the user ID from the form data
	userID := r.Form.Get("UserID")

	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	imageUrl, err := database.CopyDefaultProfilePicture(userID)
	if err != nil {
		fmt.Println("Error copying default profile picture:", err)
		http.Error(w, "Error copying default profile picture", http.StatusInternalServerError)
		return
	}

	err = database.ChangeImageUrl(db, userID, imageUrl)
	if err != nil {
		fmt.Println("Error changing image URL:", err)
		http.Error(w, "Error changing image URL", http.StatusInternalServerError)
		return
	}

	// Check the referer header
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}

	// Set cookies if not on admin page
	if filepath.Base(referer) != "admin" {
		// Get the account from the database
		acc, err := database.GetAccountbyID(db, userID)
		if err != nil {
			fmt.Println("Error getting account by username:", err)
			return
		}
		// Create a new cookie for the account
		accountCookie := &http.Cookie{
			Name: "account",
			// The value of the cookie is a string that contains the account's information separated by "|"
			Value: fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%t|%s", acc.Id, acc.Email, acc.Password, acc.Username, acc.ImageUrl, acc.IsBan, acc.IsModerator, acc.IsAdmin, acc.CreationDate),
			Path:  "/",
		}

		// Set the cookie
		http.SetCookie(w, accountCookie)
	}

	// Redirect the user to the previous page
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func ReportedPostsForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the post ID from the form data
	postID := r.Form.Get("PostID")

	db, err := database.ConnectPostDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Increment the number of views of the post
	err = database.IncrementNumberOfReportstoDB(db, postID)
	if err != nil {
		fmt.Println("Error incrementing views:", err)
		return
	}

	// Check the referer header
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}

	// Redirect to the post page
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
