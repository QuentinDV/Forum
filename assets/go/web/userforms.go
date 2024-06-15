package web

import (
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
	"strings"
)

type OtherUserProfileData struct {
	Username           string
	ImageUrl           string
	CreationDate       string
	SubscribedCategory []database.Category
	LikedPosts         []database.Post
	DisLikedPosts      []database.Post
}

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

	fmt.Println("account:", account)
	fmt.Println("Acc:", Acc)
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

	// Get the URL of the previous page from the Referer header
	previousPage := r.Header.Get("Referer")

	// Redirect to the previous page
	http.Redirect(w, r, previousPage, http.StatusSeeOther)
}

func PfpWithImageForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(4 << 20) // Set maxMemory parameter to 4MB
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

	// Redirect to the home page
	http.Redirect(w, r, "/userprofile", http.StatusSeeOther)
}

func ChangePwForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)

	oldpassword := r.Form.Get("oldPw")
	newpassword := r.Form.Get("newPw")

	// fmt.Println("oldpassword:", oldpassword)
	// fmt.Println("newpassword:", newpassword)

	// Change the image url in the database
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()
	database.ChangePassword(db, ConnectedAccount.Id, ConnectedAccount.Username, oldpassword, newpassword)

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, ConnectedAccount.Username)
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

func UserProfileForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	AccUsername := r.Form.Get("AccUsername")

	ConnectedAccount := RetrieveAccountfromCookie(r)

	if ConnectedAccount.Username == "Guest" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if ConnectedAccount.Username == AccUsername {
		http.Redirect(w, r, "/myprofile", http.StatusSeeOther)
		return
	}

	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, AccUsername)
	if err != nil {
		fmt.Println("Error getting account by username:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Get the subscribed categories of the account
	subscribedCategoriesIDs, err := database.GetSubscribedCategories(db, Acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}
	// fmt.Println("subscribedCategories:", subscribedCategoriesIDs)
	var subscribedCategories []database.Category

	for i := 1; i < len(subscribedCategoriesIDs); i++ {
		post, err := database.GetCategorybyID(db, subscribedCategoriesIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		subscribedCategories = append(subscribedCategories, post)
	}

	// Get the favorite posts of the account
	favoritePostsIDs, err := database.GetLikedPosts(db, Acc.Id)
	if err != nil {
		fmt.Println("Error getting liked posts:", err)
		return
	}
	// fmt.Println("favoritePostsIDs:", favoritePostsIDs)
	var likesPosts []database.Post

	for i := 1; i < len(favoritePostsIDs); i++ {
		post, err := database.GetPost(db, favoritePostsIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		likesPosts = append(likesPosts, post)
	}

	// Get the disliked posts of the account
	dislikedPostsIDs, err := database.GetDisLikedPosts(db, Acc.Id)
	if err != nil {
		fmt.Println("Error getting disliked posts:", err)
		return
	}
	// fmt.Println("dislikedPostsIDs:", dislikedPostsIDs)
	var dislikesPosts []database.Post

	for i := 1; i < len(dislikedPostsIDs); i++ {
		post, err := database.GetPost(db, dislikedPostsIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		dislikesPosts = append(dislikesPosts, post)
	}

	// Create a new UserProfileData struct
	userProfileData := OtherUserProfileData{
		Username:           Acc.Username,
		ImageUrl:           Acc.ImageUrl,
		CreationDate:       Acc.CreationDate,
		SubscribedCategory: subscribedCategories,
		LikedPosts:         likesPosts,
		DisLikedPosts:      dislikesPosts,
	}
	// fmt.Println("userProfileData:", userProfileData)

	// Execute the user profile template with the UserProfileData struct
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile.html"))
	tmpl.Execute(w, userProfileData)
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

func CreateCategoryForm(w http.ResponseWriter, r *http.Request) {

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve form values
	title := r.Form.Get("categoryName")
	description := r.Form.Get("description")
	imageUrl := r.Form.Get("imageUrl")
	existingTags := r.Form["existingTags[]"]
	newTags := r.Form.Get("newTags")

	// Process existing tags (if any selected)
	var tags []string
	if len(existingTags) > 0 {
		tags = existingTags
	}

	// Process new tags entered by the user
	if newTags != "" {
		newTagsSlice := strings.Split(newTags, ",")
		for _, tag := range newTagsSlice {
			trimmedTag := strings.TrimSpace(tag)
			if trimmedTag != "" {
				tags = append(tags, trimmedTag)
			}
		}
	}

	// Example: Retrieve account ID from cookie
	accountID := RetrieveAccountfromCookie(r).Id

	// Print to console (for debugging)
	fmt.Println("title:", title)
	fmt.Println("description:", description)
	fmt.Println("imageUrl:", imageUrl)
	fmt.Println("tags:", tags)
	fmt.Println("accountID:", accountID)

	// Create the category in the database
	db, err := database.ConnectCategoriesDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	err = database.CreateCategory(db, title, description, imageUrl, tags, accountID)
	if err != nil {
		fmt.Println("Error creating category:", err)
		return
	}

	Category, err := database.GetCategoryByTitle(db, title)
	if err != nil {
		fmt.Println("Error getting category by title:", err)
		return
	}

	// Add the account to the subscribed category
	err = database.AddSubscribedCategory(db, accountID, Category.CategoryID)
	if err != nil {
		fmt.Println("Error adding subscribed category:", err)
		return
	}

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
