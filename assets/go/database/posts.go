package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	Id               string
	Title            string
	Content          string
	ImageUrl         string
	Likes            int
	Dislikes         int
	View             int
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
	_, err = db.Exec("INSERT INTO posts (postID, title, content, imageUrl, likes, dislikes, view, categoryID, AccountID, creationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newID, post.Title, post.Content, post.ImageUrl, post.Likes, post.Dislikes, post.View, post.CategoryID, post.AccountID, post.CreationDate)
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
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
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
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
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
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
		FROM posts p
		JOIN categories c ON p.categoryID = c.CategoryId
		JOIN accounts a ON p.AccountID = a.id
		WHERE p.postID = ?
	`, id)
	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

// GetPostsByCategory function retrieves all posts from a specific category.
func GetPostsByCategory(db *sql.DB, categoryID string) ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
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
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
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
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
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
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
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
		SELECT p.postID, p.title, p.content, p.imageUrl, p.likes, p.dislikes, p.view, p.categoryID, c.title as categoryName, c.ImageUrl as categoryImageUrl, p.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, p.creationDate 
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
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.ImageUrl, &post.Likes, &post.Dislikes, &post.View, &post.CategoryID, &post.CategoryName, &post.CategoryImageUrl, &post.AccountID, &post.AccountUsername, &post.AccountImageUrl, &post.CreationDate)
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
