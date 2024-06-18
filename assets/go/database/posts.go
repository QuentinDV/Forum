package database

import (
	"database/sql"
	"fmt"
	"sort"
	"time"
)

type Post struct {
	PostID           string
	Title            string
	Content          string
	ImageUrl         string
	Likes            int
	Dislikes         int
	View             int
	Responses        int
	Reports          int
	CategoryID       string
	CategoryName     string
	CategoryImageUrl string
	AccountID        string
	AccountUsername  string
	AccountImageUrl  string
	CreationDate     string
}

// ConnectPostDB function creates a new connection to the SQLite database
func ConnectPostDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Creating the posts table if it does not already exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		postID TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		imageUrl TEXT,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		view INTEGER DEFAULT 0,
		responses INTEGER DEFAULT 0,
		reports INTEGER DEFAULT 0,
		categoryID TEXT NOT NULL,
		AccountID TEXT NOT NULL,
		creationDate TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InsertPost function inserts a new post into the database.
func InsertPost(db *sql.DB, post Post) error {
	// Get the last post ID
	row := db.QueryRow("SELECT postID FROM posts ORDER BY postID DESC LIMIT 1")
	var lastID string
	err := row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If there are no posts in the database, set the last ID to "0"
			lastID = "-1"
		} else {
			// If there's another error, return it
			return err
		}
	}

	// Increment the last ID
	newID := incrementID(lastID)

	// Insert the new post with the incremented ID
	_, err = db.Exec("INSERT INTO posts (postID, title, content, imageUrl, likes, dislikes, view, responses, reports, categoryID, AccountID, creationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newID, post.Title, post.Content, post.ImageUrl, post.Likes, post.Dislikes, post.View, post.Responses, post.Reports, post.CategoryID, post.AccountID, post.CreationDate)
	return err
}

// CreatePost function creates a new post in the database.
func CreatePost(db *sql.DB, title, content, imageUrl, categoryID, accountID string) error {
	post := Post{
		Title:        title,
		Content:      content,
		ImageUrl:     imageUrl,
		CategoryID:   categoryID,
		AccountID:    accountID,
		CreationDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := IncrementNumberOfPosts(db, categoryID)
	if err != nil {
		fmt.Println("Error incrementing number of posts:", err)
		return err
	}

	return InsertPost(db, post)
}

// DeletePost function deletes a post from the database.
func DeletePost(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM posts WHERE postID = ?", id)
	return err
}

// GetAllPosts function retrieves all posts from the database.
func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.responses, p.reports, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.Responses, &post.Reports, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPost function retrieves a post from the database.
func GetPost(db *sql.DB, id string) (Post, error) {
	row := db.QueryRow(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.responses, p.reports, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
		WHERE p.postID = ?
	`, id)
	var post Post
	err := row.Scan(&post.PostID, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.Responses, &post.Reports, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

// GetPostsByCategory function retrieves all posts from a specific category.
func GetPostsByCategory(db *sql.DB, categoryID string) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.responses, p.reports, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
		WHERE p.categoryID = ?
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.Responses, &post.Reports, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostsByCreator function retrieves all posts from a specific creator.
func GetPostsByCreator(db *sql.DB, AccountID string) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.responses, p.reports, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
		WHERE p.AccountID = ?
	`, AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.Responses, &post.Reports, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostsByTitle function retrieves all posts with a specific title.
func GetPostsByTitle(db *sql.DB, title string) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.responses, p.reports, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
		WHERE p.title = ?
	`, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.Responses, &post.Reports, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// AddLike function increments the number of likes of a post.
func AddLiketoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET likes = likes + 1 WHERE postID = ?", postID)
	return err
}

// RemoveLike function decrements the number of likes of a post.
func RemoveLiketoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET likes = likes - 1 WHERE postID = ?", postID)
	return err
}

// AddDislike function increments the number of dislikes of a post.
func AddDisliketoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE postID = ?", postID)
	return err
}

// RemoveDislike function decrements the number of dislikes of a post.
func RemoveDisliketoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE postID = ?", postID)
	return err
}

// IncrementView function increments the number of views of a post.
func IncrementViewtoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET view = view + 1 WHERE postID = ?", postID)
	return err
}

// IncrementNumberOfResponse function increments the number of responses of a post.
func IncrementNumberOfResponsetoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET responses = responses + 1 WHERE postID = ?", postID)
	return err
}

// DecrementNumberOfResponse function decrements the number of responses of a post.
func DecrementNumberOfResponsetoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET responses = responses - 1 WHERE postID = ?", postID)
	return err
}

// IncrementNumberOfReports function increments the number of reports of a post.
func IncrementNumberOfReportstoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET Reports = Reports + 1 WHERE postID = ?", postID)
	return err
}

// DecrementNumberOfReports function decrements the number of reports of a post.
func DecrementNumberOfReportstoDB(db *sql.DB, postID string) error {
	_, err := db.Exec("UPDATE posts SET Reports = Reports - 1 WHERE postID = ?", postID)
	return err
}

// Generate a new post ID
func GenerateNewPostID(db *sql.DB) string {
	// Get the last post ID
	row := db.QueryRow("SELECT postID FROM posts ORDER BY postID DESC LIMIT 1")
	var lastID string
	err := row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If there are no posts in the database, set the last ID to "0"
			lastID = "-1"
		} else {
			// If there's another error, return an empty string (handle appropriately in the caller)
			return ""
		}
	}

	// Increment the last ID
	newID := incrementID(lastID)
	return newID
}

// SortPostsByCreationDate function sorts posts by creation date in descending order.
func SortPostsByDateDescending(posts []Post) ([]Post, error) {
	// Parse the CreationDate field to time.Time for sorting
	parsedPosts := make([]Post, len(posts))
	for i, post := range posts {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", post.CreationDate)
		if err != nil {
			return nil, err
		}
		parsedPosts[i] = post
		parsedPosts[i].CreationDate = parsedTime.Format(time.RFC3339) // Use a standard time format for sorting
	}

	// Sort the posts by CreationDate
	sort.Slice(parsedPosts, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, parsedPosts[i].CreationDate)
		timeJ, _ := time.Parse(time.RFC3339, parsedPosts[j].CreationDate)
		return timeI.After(timeJ) // Sort in descending order
	})

	return parsedPosts, nil
}

// SortPostsByDateAscending sorts a slice of posts by their creation date in ascending order (oldest first).
func SortPostsByDateAscending(posts []Post) ([]Post, error) {
	// Parse the CreationDate field to time.Time for sorting
	parsedPosts := make([]Post, len(posts))
	for i, post := range posts {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", post.CreationDate)
		if err != nil {
			return nil, err
		}
		parsedPosts[i] = post
		parsedPosts[i].CreationDate = parsedTime.Format(time.RFC3339) // Use a standard time format for sorting
	}

	// Sort the posts by CreationDate
	sort.Slice(parsedPosts, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, parsedPosts[i].CreationDate)
		timeJ, _ := time.Parse(time.RFC3339, parsedPosts[j].CreationDate)
		return timeI.Before(timeJ) // Sort in ascending order
	})

	return parsedPosts, nil
}

// Function to sort posts by the number of likes
func DescendingPostsSortingByLikes(posts []Post) []Post {
	// Sort the posts by the number of likes in descending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Likes > posts[j].Likes
	})

	return posts
}

// Function to sort posts by the number of likes in ascending order
func AscendingPostsSortingByLikes(posts []Post) []Post {
	// Sort the posts by the number of likes in ascending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Likes < posts[j].Likes
	})

	return posts
}

// FilterPostsByCategory filters posts and keeps those with a specific category name.
func FilterPostsByCategory(posts []Post, categoryName string) []Post {
	filteredPosts := []Post{}

	for _, post := range posts {
		if post.CategoryName == categoryName {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

// SortPostsByViews sorts a slice of posts by their views in ascending order.
func SortPostsByViewsAscending(posts []Post) []Post {
	// Sort the posts by View
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].View < posts[j].View
	})

	return posts
}

// SortPostsByViewsDescending sorts a slice of posts by their views in descending order.
func SortPostsByViewsDescending(posts []Post) []Post {
	// Sort the posts by View in descending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].View > posts[j].View
	})

	return posts
}

// SortPostsByResponses sorts a slice of posts by their responses in ascending order.
func SortPostsByResponsesAscending(posts []Post) []Post {
	// Sort the posts by Responses
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Responses < posts[j].Responses
	})

	return posts
}

// SortPostsByResponsesDescending sorts a slice of posts by their responses in descending order.
func SortPostsByResponsesDescending(posts []Post) []Post {
	// Sort the posts by Responses in descending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Responses > posts[j].Responses
	})

	return posts
}
