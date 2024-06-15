package web

import (
	"fmt"
	"forum/assets/go/database"
	"net/http"
	"strings"
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

// CreatePostForm function handles the form submission for creating a new post
func CreatePostForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Retrieve form values
	title := r.Form.Get("postTitle")
	content := r.Form.Get("postContent")
	imageUrl := r.Form.Get("imageUrl")
	categoryName := r.Form.Get("categoryName")

	// Example: Retrieve account ID from cookie
	accountID := RetrieveAccountfromCookie(r).Id

	// Print to console (for debugging)
	fmt.Println("title:", title)
	fmt.Println("content:", content)
	fmt.Println("imageUrl:", imageUrl)
	fmt.Println("categoryName:", categoryName)
	fmt.Println("accountID:", accountID)

	// Create the post in the database
	db, err := database.ConnectPostDB("db/database.db")
	if err != nil {
		return
	}
	defer db.Close()

	category, err := database.GetCategoryByTitle(db, categoryName)
	if err != nil {
		fmt.Println("Error getting category by title:", err)
		return
	}

	err = database.CreatePost(db, title, content, imageUrl, category.CategoryID, accountID)
	if err != nil {
		fmt.Println("Error creating post:", err)
		return
	}

	// Redirect to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
