package web

import (
	"fmt"
	"forum/assets/go/database"
	"log"
	"net/http"
	"strings"
	"time"
)

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

// SubscribeCategoryForm handles subscribing and unsubscribing from a category.
func SubscribeCategoryForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve form values
	categoryID := r.Form.Get("categoryID")

	// Retrieve account ID from cookie
	accountID := RetrieveAccountfromCookie(r).Id

	// Connect to the database
	db, err := database.ConnectCategoriesDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Check if the category is already subscribed
	IsThisCategorySubscribed := database.IsThisCategorySubscribed(db, accountID, categoryID)
	if IsThisCategorySubscribed {
		// Remove the subscribed category if already subscribed
		err = database.RemoveSubscribedCategory(db, accountID, categoryID)
		if err != nil {
			fmt.Println("Error removing subscribed category:", err)
			http.Error(w, "Error removing subscribed category", http.StatusInternalServerError)
			return
		}
	} else {
		// Add the account to the subscribed category
		err = database.AddSubscribedCategory(db, accountID, categoryID)
		if err != nil {
			fmt.Println("Error adding subscribed category:", err)
			http.Error(w, "Error adding subscribed category", http.StatusInternalServerError)
			return
		}
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func LikeForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	PostID := r.Form.Get("LikeID")

	if ConnectedAccount.Username == "Guest" {
		// Redirect the user back to the previous page;
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/" // Fallback URL if Referer header is not set
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

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

	if ConnectedAccount.Username == "Guest" {
		// Redirect the user back to the previous page
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/" // Fallback URL if Referer header is not set
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

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

// CreatePostForm handles the form submission for creating a new post
func CreatePostForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(4 << 20) // Set maxMemory parameter to 4MB
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve the connected account from the cookie
	ConnectedAccount := RetrieveAccountfromCookie(r)
	if (ConnectedAccount == database.Account{}) {
		http.Error(w, "Unable to retrieve account from cookie", http.StatusUnauthorized)
		return
	}

	// Retrieve form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryName := r.FormValue("categoryName")

	// Create the post in the database
	db, err := database.ConnectPostDB("db/database.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	category, err := database.GetCategoryByTitle(db, categoryName)
	if err != nil {
		http.Error(w, "Error retrieving category", http.StatusInternalServerError)
		return
	}

	// Generate a new post ID
	postID := database.GenerateNewPostID(db)

	// Get the profile picture file from the form data
	file, _, err := r.FormFile("postimage")
	var imageUrl string
	if err == nil {
		// File is present, save it
		defer file.Close()
		filePath := fmt.Sprintf("./assets/img/post/%s.png", postID)
		imageUrl = fmt.Sprintf("./assets/img/post/%s.png", postID)
		err = database.SaveFile(filePath, file)
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
	} else {
		// File is not present, set a default or empty URL
		imageUrl = ""
	}

	// Insert the new post into the database
	post := database.Post{
		PostID:       postID,
		Title:        title,
		Content:      content,
		ImageUrl:     imageUrl,
		CategoryID:   category.CategoryID,
		AccountID:    ConnectedAccount.Id,
		CreationDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = database.InsertPost(db, post)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func SavePostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	PostID := r.Form.Get("PostID")

	db, err := database.ConnectUserDataDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	IsThisPostSaved := database.IsThisPostSaved(db, ConnectedAccount.Id, PostID)
	if IsThisPostSaved {
		err = database.RemoveSavedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error removing saved post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	} else {
		err = database.AddSavedPost(db, ConnectedAccount.Id, PostID)
		if err != nil {
			fmt.Println("Error adding saved post:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// CreateCommentForm handles the form submission for creating a new comment
func CreateCommentForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(4 << 20) // Set maxMemory parameter to 4MB
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve the connected account from the cookie
	ConnectedAccount := RetrieveAccountfromCookie(r)
	if (ConnectedAccount == database.Account{}) {
		http.Error(w, "Unable to retrieve account from cookie", http.StatusUnauthorized)
		return
	}

	// Retrieve form values
	content := r.FormValue("content")
	postID := r.FormValue("PostID")

	// Create the comment in the database
	db, err := database.ConnectCommentsDB("db/database.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get the profile picture file from the form data
	file, _, err := r.FormFile("postimage")
	var imageUrl string
	if err == nil {
		// File is present, save it
		defer file.Close()

		commentID, err := database.GenerateNewCommentID(db)
		if err != nil {
			http.Error(w, "Error generating new comment ID", http.StatusInternalServerError)
			return
		}

		filePath := fmt.Sprintf("./assets/img/comment/%s.png", commentID)
		imageUrl = filePath
		err = database.SaveFile(filePath, file)
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
	} else {
		// File is not present, set a default or empty URL
		imageUrl = ""
	}

	// Insert the new comment into the database
	comment := database.Comment{
		PostID:       postID,
		Content:      content,
		ImageUrl:     imageUrl,
		AccountID:    ConnectedAccount.Id,
		CreationDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = database.InsertComment(db, comment)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func LikeCommentForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	CommentID := r.Form.Get("CommentID")

	if ConnectedAccount.Username == "Guest" {
		// Redirect the user back to the previous page
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/" // Fallback URL if Referer header is not set
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	db, err := database.ConnectUserDataDB("db/database.db")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	// Get the comment's current likes and dislikes
	Comment, err := database.GetComment(db, CommentID)
	if err != nil {
		log.Println("Error fetching comment:", err)
		http.Error(w, "Error fetching comment", http.StatusInternalServerError)
		return
	}
	log.Println("Current likes and dislikes for comment:", Comment.Likes, Comment.Dislikes)

	// Check if the comment is already liked by the connected account
	IsThisCommentLiked := database.IsThisCommentLiked(db, ConnectedAccount.Id, CommentID)

	log.Println("IsThisCommentLiked:", IsThisCommentLiked)

	// Toggle like/unlike based on current like status
	if IsThisCommentLiked {
		err = database.RemoveLikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			log.Println("Error removing liked comment:", err)
			http.Error(w, "Error removing liked comment", http.StatusInternalServerError)
			return
		}
		log.Println("Comment unliked successfully")
	} else {
		err = database.RemoveDislikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			log.Println("Error removing disliked comment:", err)
			http.Error(w, "Error removing disliked comment", http.StatusInternalServerError)
			return
		}

		err = database.AddLikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			log.Println("Error adding liked comment:", err)
			http.Error(w, "Error adding liked comment", http.StatusInternalServerError)
			return
		}
		log.Println("Comment liked successfully")
	}

	// Redirect the user back to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeCommentForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	ConnectedAccount := RetrieveAccountfromCookie(r)
	CommentID := r.Form.Get("CommentID")

	if ConnectedAccount.Username == "Guest" {
		// Redirect the user back to the previous page
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/" // Fallback URL if Referer header is not set
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	db, err := database.ConnectUserDataDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	IsThisCommentDisliked := database.IsThisCommentDisliked(db, ConnectedAccount.Id, CommentID)
	if IsThisCommentDisliked {
		err = database.RemoveDislikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			fmt.Println("Error removing disliked comment:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	} else {
		err = database.RemoveLikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			fmt.Println("Error removing liked comment:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}

		err = database.AddDislikedComment(db, ConnectedAccount.Id, CommentID)
		if err != nil {
			fmt.Println("Error adding disliked comment:", err)
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DeleteCommentForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	CommentID := r.Form.Get("CommentID")

	db, err := database.ConnectCommentsDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	err = database.DeleteComment(db, CommentID)
	if err != nil {
		fmt.Println("Error deleting comment:", err)
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the previous page
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/" // Fallback URL if Referer header is not set
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// SortPosts trie les posts en fonction de la méthode de tri sélectionnée.
func SortPosts(posts []database.Post, sortedBy string) ([]database.Post, error) {
	var err error

	switch sortedBy {
	case "By Date Descending":
		posts, err = database.SortPostsByDateDescending(posts)
	case "By Date Ascending":
		posts, err = database.SortPostsByDateAscending(posts)
	case "By Likes Descending":
		posts = database.DescendingPostsSortingByLikes(posts)
	case "By Likes Ascending":
		posts = database.AscendingPostsSortingByLikes(posts)
	case "By Views Ascending":
		posts = database.SortPostsByViewsAscending(posts)
	case "By Views Descending":
		posts = database.SortPostsByViewsDescending(posts)
	case "By Responses Ascending":
		posts = database.SortPostsByResponsesAscending(posts)
	case "By Responses Descending":
		posts = database.SortPostsByResponsesDescending(posts)
	default:
		return nil, fmt.Errorf("invalid sorting method: %s", sortedBy)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// SortingHomePostsForm handles the form submission for sorting the home posts
func SortingHomePostsForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve the sorting method from the form
	SortedBy = r.Form.Get("sortingMethod")
	CategoryName = r.Form.Get("categoryName")

	fmt.Println("SortedBy:", SortedBy)
	fmt.Println("CategoryName:", CategoryName)

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}

// ResetHOmeSorting resets the sorting method to the default value
func ResetHomeSortingForm(w http.ResponseWriter, r *http.Request) {
	SortedBy = "By Date Descending"
	CategoryName = ""

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
