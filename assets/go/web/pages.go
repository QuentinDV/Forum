package web

// Importing necessary packages
import (
	"database/sql"
	"fmt"
	"forum/assets/go/database"
	"html/template"
	"net/http"
	"strings"
)

// HomeData struct represents the data needed to render the home page
type HomeData struct {
	Username            string
	ImageUrl            string
	SortedBy            string
	FavoritesCategories []database.Category
	AllCategories       []database.Category
	AllPosts            []database.Post
	RecentPosts         []database.Post
	ExistingTags        []string
}

type UserProfile struct {
	Username                     string
	ImageUrl                     string
	CreationDate                 string
	Email                        string
	ID                           string
	IsSameAccount                bool
	IsAdmin                      bool
	NumberofSubscribedCategories int
	MyPosts                      []database.Post
	LikedPosts                   []database.Post
	DislikedPosts                []database.Post
	MyComments                   []database.Comment
	LikedComments                []database.Comment
	DislikedComments             []database.Comment
	SavedPosts                   []database.Post
}

type CommentWithPermissions struct {
	database.Comment
	IsAdmin     bool
	IsModerator bool
	IsSameUser  bool
}

type PostData struct {
	IsAdmin     bool
	IsModerator bool
	IsSameUser  bool
	IsSaved     bool
	IsGuest     bool
	Post        database.Post
	Comments    []CommentWithPermissions
}

// CategoryData struct represents the data needed to render the category page
type CategoryData struct {
	Category      database.Category
	Posts         []database.Post
	CategoryTags  []string
	ExistingTags  []string
	AllUsernames  []string
	IsSubscribed  bool
	IsAdmin       bool
	IsSameAccount bool
	Username      string
}

type CreateCategoryData struct {
	ExistingTags []string
}

type CreatePostData struct {
	CategoryTitles []string
}

var SortedBy = "By Date Descending"
var CategoryName string

// Home is the main page of the forum.
func Home(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Open the database
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	defer db.Close()

	// Get all categories from the database
	allCategories, err := database.SortBySubsriber(db)
	if err != nil {
		fmt.Println("Error getting all categories:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Get the favorite categories of the account
	favoriteCategoriesIDs, err := database.GetSubscribedCategories(db, ConnectedAccount.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	var favoriteCategories []database.Category

	for i := 1; i < len(favoriteCategoriesIDs); i++ {
		post, err := database.GetCategorybyID(db, favoriteCategoriesIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		favoriteCategories = append(favoriteCategories, post)
	}

	// Get all tags from the database
	allTag, err := database.GetAllTags(db)
	if err != nil {
		fmt.Println("Error getting all tags:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Get the recent posts
	allPosts, err := database.GetAllPosts(db)
	if err != nil {
		fmt.Println("Error getting recent posts:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	if SortedBy == "By Categories" {
		allPosts = database.FilterPostsByCategory(allPosts, CategoryName)
	} else {
		allPosts, err = SortPosts(allPosts, SortedBy)
		if err != nil {
			fmt.Println("Error sorting posts:", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	}

	// Create a new HomeData struct
	HomeData := HomeData{
		Username:            ConnectedAccount.Username,
		ImageUrl:            ConnectedAccount.ImageUrl,
		SortedBy:            SortedBy,
		FavoritesCategories: favoriteCategories,
		AllCategories:       allCategories,
		AllPosts:            allPosts,
		ExistingTags:        allTag,
	}

	// Execute the home template with the HomeData struct
	tmpl := template.Must(template.ParseFiles("assets/html/home.html"))
	err = tmpl.Execute(w, HomeData)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
}

// LogIn page of the forum.
func LogIn(w http.ResponseWriter, r *http.Request) {
	// Serve the login page
	http.ServeFile(w, r, "assets/html/login.html")
}

// Admin page of the forum.
func Admin(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Open the database
	db, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	allAcc, err := database.GetAllAccounts(db)
	if err != nil {
		fmt.Println("Error getting all accounts:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	if !ConnectedAcc.IsAdmin || ConnectedAcc.Username == "Guest" {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Serve the admin page
	tmpl := template.Must(template.ParseFiles("assets/html/admin.html"))
	tmpl.Execute(w, allAcc)
}

// Handler to render the ReportePosts page
func ReportePosts(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user is the same as the connected account
	if !ConnectedAccount.IsAdmin || ConnectedAccount.Username == "Guest" {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Open the database
	db, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// get all posts of the account
	Reportedposts, err := database.GetReportedPosts(db)
	if err != nil {
		fmt.Println("Error getting posts by creator:", err)
		return
	}

	Reportedposts = database.SortPostsByReportsDescending(Reportedposts)

	// Serve the admin page
	tmpl := template.Must(template.ParseFiles("assets/html/reportedposts.html"))
	tmpl.Execute(w, Reportedposts)
}

// Handler to render the create post page
func CreatePostHome(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Open the database
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	defer db.Close()

	// Get all category titles from the database
	allCategories, err := database.GetAllCategories(db)
	if err != nil {
		fmt.Println("Error getting all categories:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Extract category titles
	var categoryTitles []string
	for _, category := range allCategories {
		categoryTitles = append(categoryTitles, category.Title)
	}

	// Prepare data for the template
	data := CreatePostData{
		CategoryTitles: categoryTitles,
	}

	// Render the create category page
	tmpl := template.Must(template.ParseFiles("assets/html/creation/categorycreation.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
}

// 404 page of the forum.
func NotFound(w http.ResponseWriter, r *http.Request) {
	// Serve the 404 page
	http.ServeFile(w, r, "assets/html/404.html")
}

// CategoryPageHandler handles the category page requests
func CategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extrait l'ID de la catégorie de l'URL
	CategoryID := r.URL.Path[len("/category/"):]

	// Vérifiez que l'ID de la catégorie n'est pas vide et ne commence pas par "assets/img/pfp/"
	if CategoryID == "" || strings.HasPrefix(CategoryID, "assets/img/pfp/") {
		http.NotFound(w, r)
		return
	}

	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	defer db.Close()

	// Récupère la catégorie de la base de données
	category, err := database.GetCategorybyID(db, CategoryID)
	if err != nil {
		fmt.Println("Error getting category by ID:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Vérifie si l'utilisateur est abonné à la catégorie
	isSubscribed := database.IsThisCategorySubscribed(db, RetrieveAccountfromCookie(r).Id, category.CategoryID)

	// Récupère les posts associés à la catégorie
	posts, err := database.GetPostsByCategory(db, category.CategoryID)
	if err != nil {
		fmt.Println("Error getting posts by category ID:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Get all tags from the database
	allTag, err := database.GetAllTags(db)
	if err != nil {
		fmt.Println("Error getting all tags:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	AllUsernames, err := database.GetAllUsernames(db)
	if err != nil {
		fmt.Println("Error getting all usernames:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	CategoryData := struct {
		Category      database.Category
		Posts         []database.Post
		IsSubscribed  bool
		IsAdmin       bool
		IsSameAccount bool
		Username      string
		ExistingTags  []string
		CategoryTags  []string
		AllUsernames  []string
	}{
		Category:      category,
		Posts:         posts,
		ExistingTags:  allTag,
		CategoryTags:  category.Tags,
		AllUsernames:  AllUsernames,
		IsSubscribed:  isSubscribed,
		IsAdmin:       ConnectedAccount.IsAdmin,
		IsSameAccount: category.AccountID == ConnectedAccount.Id,
		Username:      ConnectedAccount.Username,
	}

	// Execute the user profile template with the CategoryData struct
	tmpl := template.Must(template.ParseFiles("assets/html/category.html"))
	err = tmpl.Execute(w, CategoryData)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
		return
	}
}

// PostPageHandler handles the post page.
func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extract the post ID from the URL
	PostID := r.URL.Path[len("/post/"):]

	// Verify that the post ID is not empty and does not start with "assets/img/pfp/"
	if PostID == "" || strings.HasPrefix(PostID, "assets/img/pfp/") {
		http.NotFound(w, r)
		return
	}

	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	defer db.Close()

	// Get the post from the database
	post, err := database.GetPost(db, PostID)
	if err != nil {
		fmt.Println("Error getting post by ID:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Get all comments of the post
	comments, err := database.GetCommentsByPost(db, PostID)
	if err != nil {
		fmt.Println("Error getting comments by post ID:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Retrieve connected account
	connectedAccount := RetrieveAccountfromCookie(r)

	// Parcourez chaque commentaire et ajoutez les informations de permission
	var commentsWithPermissions []CommentWithPermissions
	for _, comment := range comments {
		commentsWithPermissions = append(commentsWithPermissions, CommentWithPermissions{
			Comment:     comment,
			IsAdmin:     connectedAccount.IsAdmin,
			IsModerator: connectedAccount.IsModerator,
			IsSameUser:  connectedAccount.Id == comment.AccountID,
		})
	}

	data := PostData{
		IsAdmin:     connectedAccount.IsAdmin,
		IsModerator: connectedAccount.IsModerator,
		IsSameUser:  post.AccountID == connectedAccount.Id,
		IsSaved:     database.IsThisPostSaved(db, connectedAccount.Id, post.PostID),
		IsGuest:     connectedAccount.Username == "Guest",
		Post:        post,
		Comments:    commentsWithPermissions,
	}

	// Execute the user profile template with the PostPageHandler struct
	tmpl := template.Must(template.ParseFiles("assets/html/post.html"))
	tmpl.Execute(w, data)
}

// UserProfileHandler handles the user profile page and its subpages.
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/user/"), "/")

	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}

	username := parts[0]
	page := ""
	if len(parts) > 1 {
		page = parts[1]
	}

	// Check if the username is valid
	if username == "" || strings.HasPrefix(username, "assets") || username == "Guest" {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Open the database connection
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	if username == ConnectedAcc.Username {
		// Route to the appropriate handler based on the page
		switch page {
		case "":
			handleProfileMainPage(w, r, db, ConnectedAcc)
		case "liked":
			handleLikedPostsPage(w, r, db, ConnectedAcc)
		case "disliked":
			handleDislikedPostsPage(w, r, db, ConnectedAcc)
		case "comments":
			handleCommentsPage(w, r, db, ConnectedAcc)
		case "savedposts":
			handleSavedPostsPage(w, r, db, ConnectedAcc)
		case "account":
			handleAccountPage(w, r, db, ConnectedAcc)
		default:
			http.NotFound(w, r)
		}
		return
	}

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, username)
	if err != nil {
		fmt.Println("Error getting account by username:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Route to the appropriate handler based on the page
	switch page {
	case "":
		handleProfileMainPage(w, r, db, Acc)
	case "liked":
		handleLikedPostsPage(w, r, db, Acc)
	case "disliked":
		handleDislikedPostsPage(w, r, db, Acc)
	case "comments":
		handleCommentsPage(w, r, db, Acc)
	case "savedposts":
		handleSavedPostsPage(w, r, db, Acc)
	case "account":
		handleAccountPage(w, r, db, Acc)
	default:
		http.NotFound(w, r)
	}
}

func MyProfile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Open the database
	db, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	if ConnectedAcc.Username == "" || ConnectedAcc.Username == "Guest" {
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Determine which page to display based on the URL path
	path := r.URL.Path
	switch path {
	case "/myprofile":
		handleProfileMainPage(w, r, db, ConnectedAcc)
	case "/userprofile/liked":
		handleLikedPostsPage(w, r, db, ConnectedAcc)
	case "/userprofile/disliked":
		handleDislikedPostsPage(w, r, db, ConnectedAcc)
	case "/userprofile/comments":
		handleCommentsPage(w, r, db, ConnectedAcc)
	case "/userprofile/savedposts":
		handleSavedPostsPage(w, r, db, ConnectedAcc)
	case "/userprofile/account":
		handleAccountPage(w, r, db, ConnectedAcc)
	default:
		http.NotFound(w, r)
	}
}

func handleProfileMainPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// get all posts of the account
	AccountPosts, err := database.GetPostsByCreator(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting posts by creator:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	}

	// Create a new UserProfile struct
	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
		MyPosts:                      AccountPosts,
	}

	// Serve the main profile page template
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/main.html"))
	tmpl.Execute(w, data)
}

func handleLikedPostsPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the favorite posts of the account
	favoritePostsIDs, err := database.GetLikedPosts(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting liked posts:", err)
		return
	}
	var likesPosts []database.Post

	for i := 1; i < len(favoritePostsIDs); i++ {
		post, err := database.GetPost(db, favoritePostsIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		likesPosts = append(likesPosts, post)
	}

	// Get the favorite comments of the account
	likedCommentsIDs, err := database.GetLikedComments(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting liked comments:", err)
		return
	}
	var likesComments []database.Comment

	for i := 1; i < len(likedCommentsIDs); i++ {
		comment, err := database.GetComment(db, likedCommentsIDs[i])
		if err != nil {
			fmt.Println("Error getting comment:", err)
			return
		}
		likesComments = append(likesComments, comment)
	}
	fmt.Println(likesComments)

	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	}

	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
		LikedPosts:                   likesPosts,
		LikedComments:                likesComments,
	}

	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/liked.html"))
	tmpl.Execute(w, data)
}

func handleDislikedPostsPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the disliked posts of the account
	dislikedPostsIDs, err := database.GetDisLikedPosts(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting disliked posts:", err)
		return
	}
	var dislikesPosts []database.Post

	for i := 1; i < len(dislikedPostsIDs); i++ {
		post, err := database.GetPost(db, dislikedPostsIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		dislikesPosts = append(dislikesPosts, post)
	}

	// Get the disliked comments of the account
	dislikedCommentsIDs, err := database.GetDislikedComments(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting disliked comments:", err)
		return
	}
	var dislikesComments []database.Comment

	for i := 1; i < len(dislikedCommentsIDs); i++ {
		comment, err := database.GetComment(db, dislikedCommentsIDs[i])
		if err != nil {
			fmt.Println("Error getting comment:", err)
			return
		}
		dislikesComments = append(dislikesComments, comment)
	}

	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	}

	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
		DislikedPosts:                dislikesPosts,
		DislikedComments:             dislikesComments,
	}

	// Serve the disliked posts page template
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/disliked.html"))
	tmpl.Execute(w, data)
}

func handleCommentsPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// get all posts of the account
	comments, err := database.GetCommentsByAccount(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting posts by creator:", err)
		return
	}
	fmt.Println(comments)

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	}

	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
		MyComments:                   comments,
	}

	// Serve the comments page template
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/comments.html"))
	tmpl.Execute(w, data)
}

func handleSavedPostsPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Load and prepare data specific to the saved posts page
	SavedPostsIDs, err := database.GetSavedPosts(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting disliked posts:", err)
		return
	}
	var SavedPosts []database.Post

	for i := 1; i < len(SavedPostsIDs); i++ {
		post, err := database.GetPost(db, SavedPostsIDs[i])
		if err != nil {
			fmt.Println("Error getting post:", err)
			return
		}
		SavedPosts = append(SavedPosts, post)
	}

	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	}

	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
		SavedPosts:                   SavedPosts,
	}

	// Serve the saved posts page template
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/saved.html"))
	tmpl.Execute(w, data)
}

func handleAccountPage(w http.ResponseWriter, r *http.Request, db *sql.DB, acc database.Account) {
	// Retrieve the account from cookies
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Check if the ConnectedAccount is nil or not valid
	if (ConnectedAccount == database.Account{}) || ConnectedAccount.Id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Load and prepare data specific to the account settings page
	// Get the favorite categories of the account
	NumberofSubscribedCategories, err := database.GetSubscribedCategories(db, acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}

	// Retrieve connected account
	ConnectedAcc := RetrieveAccountfromCookie(r)

	// Check if the user is the same as the connected account
	isSameAccount := false
	if acc.Username == ConnectedAcc.Username {
		isSameAccount = true
	} else {
		http.Redirect(w, r, "../../notfound", http.StatusSeeOther)
		return
	}

	data := UserProfile{
		Username:                     acc.Username,
		ImageUrl:                     acc.ImageUrl,
		CreationDate:                 acc.CreationDate,
		ID:                           acc.Id,
		IsSameAccount:                isSameAccount,
		IsAdmin:                      acc.IsAdmin,
		Email:                        acc.Email,
		NumberofSubscribedCategories: len(NumberofSubscribedCategories) - 1,
	}
	// Serve the account settings page template
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile/account.html"))
	tmpl.Execute(w, data)
}
