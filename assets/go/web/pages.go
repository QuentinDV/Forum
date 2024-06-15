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
	ConnectedAccount    database.Account
	FavoritesCategories []database.Category
	AllCategories       []database.Category
	AllPosts            []database.Post
	RecentPosts         []database.Post
	TopPosts            []database.Post
}

// CategoryData struct represents the data needed to render the category page
type CategoryData struct {
	Category database.Category
	Posts    []database.Post
}

// Home is the main page of the forum.
func Home(w http.ResponseWriter, r *http.Request) {
	ConnectedAccount := RetrieveAccountfromCookie(r)

	// Open the database
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	defer db.Close()

	// Get all categories from the database
	allCategories, err := database.GetAllCategories(db)
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

	// Get the recent posts
	allPosts, err := database.GetAllPosts(db)
	if err != nil {
		fmt.Println("Error getting recent posts:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Get the recent posts
	// recentPosts, err := database.GetRecentPosts(db)
	// if err != nil {
	// 	fmt.Println("Error getting recent posts:", err)
	// 	http.Redirect(w, r, "/error", http.StatusSeeOther)
	// 	return
	// }

	// Get the top posts
	// topPosts, err := database.GetTopPosts(db)
	// if err != nil {
	// 	fmt.Println("Error getting top posts:", err)
	// 	http.Redirect(w, r, "/error", http.StatusSeeOther)
	// 	return
	// }

	// Create a new HomeData struct
	HomeData := HomeData{
		ConnectedAccount:    ConnectedAccount,
		FavoritesCategories: favoriteCategories,
		AllCategories:       allCategories,
		AllPosts:            allPosts,
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

// Categories page of the forum.
func Categories(w http.ResponseWriter, r *http.Request) {
	// Serve the categories page
	http.ServeFile(w, r, "assets/html/categories.html")
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

	// Serve the admin page
	tmpl := template.Must(template.ParseFiles("assets/html/admin.html"))
	tmpl.Execute(w, allAcc)
}

// Categories page of the forum.
func UserProfile(w http.ResponseWriter, r *http.Request) {
	// Serve the user profile page
	http.ServeFile(w, r, "assets/html/myprofile.html")
}

// Profile page of the forum.
func OtherUserProfile(w http.ResponseWriter, r *http.Request) {
	// Serve the other user profile page
	http.ServeFile(w, r, "assets/html/userprofile.html")
}

// 404 page of the forum.
func NotFound(w http.ResponseWriter, r *http.Request) {
	// Serve the 404 page
	http.ServeFile(w, r, "assets/html/404.html")
}

// CategoryPageHandler handles the category page requests
func CategoryPageHandler(w http.ResponseWriter, r *http.Request) {
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

	// Récupère les posts associés à la catégorie
	posts, err := database.GetPostsByCategory(db, category.CategoryID)
	if err != nil {
		fmt.Println("Error getting posts by category ID:", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	CategoryData := struct {
		Category database.Category
		Posts    []database.Post
	}{
		Category: category,
		Posts:    posts,
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
	// Extrait l'ID du post de l'URL
	PostID := r.URL.Path[len("/post/"):]
	fmt.Println("PostID", PostID)

	// Vérifiez que l'ID du post n'est pas vide et ne commence pas par "assets/img/pfp/"
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

	// Récupère le post de la base de données
	post, err := database.GetPost(db, PostID)
	if err != nil {
		fmt.Println("Error getting post by ID:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Execute the user profile template with the PostPageHandler struct
	tmpl := template.Must(template.ParseFiles("assets/html/post.html"))
	tmpl.Execute(w, post)
}

// UserProfileHandler handles the user profile page.
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extrait le username de l'URL
	AccUsername := r.URL.Path[len("/user/"):]

	// Vérifiez que le nom d'utilisateur n'est pas vide et ne commence pas par "assets/img/pfp/"
	if AccUsername == "" || strings.HasPrefix(AccUsername, "assets/img/pfp/") {
		http.NotFound(w, r)
		return
	}

	// Ici, ajoutez la logique pour récupérer les données de l'utilisateur à partir de votre base de données
	// userData := getUserDataFromDB(username)
	db, err := database.ConnectUserDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	// Get the account from the database
	Acc, err := database.GetAccountByUsername(db, AccUsername)
	if err != nil {
		fmt.Println("Error getting account by username zrzz:", err)
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	// Get the subscribed categories of the account
	subscribedCategoriesIDs, err := database.GetSubscribedCategories(db, Acc.Id)
	if err != nil {
		fmt.Println("Error getting subscribed categories:", err)
		return
	}
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

	// Execute the user profile template with the UserProfileData struct
	tmpl := template.Must(template.ParseFiles("assets/html/userprofile.html"))
	tmpl.Execute(w, userProfileData)

}
