package database

import (
	"database/sql"
	"time"
)

type Comment struct {
	PostID          string
	CommentID       string
	Content         string
	ImageUrl        string
	Likes           int
	Dislikes        int
	AccountID       string
	AccountUsername string
	AccountImageUrl string
	CreationDate    string
}

// ConnectCommentsDB function creates a new connection to the SQLite database
func ConnectCommentsDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Creating the comments table if it does not already exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		postID TEXT NOT NULL,
		commentID TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		imageUrl TEXT,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		AccountID TEXT NOT NULL,
		creationDate TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InsertComment function inserts a new comment into the database.
func InsertComment(db *sql.DB, comment Comment) error {
	// Get the last comment ID
	row := db.QueryRow("SELECT commentID FROM comments ORDER BY commentID DESC LIMIT 1")
	var lastID string
	err := row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If there are no comments in the database, set the last ID to "0"
			lastID = "-1"
		} else {
			// If there's another error, return it
			return err
		}
	}
	// Increment the last ID
	newID := incrementID(lastID)

	// Insert the new comment into the database
	_, err = db.Exec("INSERT INTO comments (postID, commentID, content, imageUrl, likes, dislikes, AccountID, creationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", comment.PostID, newID, comment.Content, comment.ImageUrl, comment.Likes, comment.Dislikes, comment.AccountID, comment.CreationDate)
	if err != nil {
		return err
	}
	return nil
}

// CreateComment function creates a new comment in the database.
func CreateComment(db *sql.DB, postID, content, imageUrl string, Acc Account) error {
	comment := Comment{
		PostID:       postID,
		Content:      content,
		ImageUrl:     imageUrl,
		AccountID:    Acc.Id,
		CreationDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := InsertComment(db, comment)
	if err != nil {
		return err
	}
	return nil
}

// GetAllComments function returns all the comments from the database.
func GetAllComments(db *sql.DB, postID string) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT c.postID, c.commentID, c.content, c.imageUrl, c.likes, c.dislikes, c.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, c.creationDate
		FROM comments c
		JOIN accounts a ON c.AccountID = a.id
		WHERE c.postID = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.PostID, &comment.CommentID, &comment.Content, &comment.ImageUrl, &comment.Likes, &comment.Dislikes, &comment.AccountID, &comment.AccountUsername, &comment.AccountImageUrl, &comment.CreationDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetComment function returns a comment from the database.
func GetComment(db *sql.DB, commentID string) (Comment, error) {
	row := db.QueryRow(`
		SELECT c.postID, c.commentID, c.content, c.imageUrl, c.likes, c.dislikes, c.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, c.creationDate
		FROM comments c
		JOIN accounts a ON c.AccountID = a.id
		WHERE c.commentID = ?
	`, commentID)
	var comment Comment
	err := row.Scan(&comment.PostID, &comment.CommentID, &comment.Content, &comment.ImageUrl, &comment.Likes, &comment.Dislikes, &comment.AccountID, &comment.AccountUsername, &comment.AccountImageUrl, &comment.CreationDate)
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

// GetCommentsByAccount function returns all the comments from a specific account.
func GetCommentsByAccount(db *sql.DB, accountID string) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT c.postID, c.commentID, c.content, c.imageUrl, c.likes, c.dislikes, c.AccountID, a.username as accountUsername, a.ImageUrl as accountImageUrl, c.creationDate
		FROM comments c
		JOIN accounts a ON c.AccountID = a.id
		WHERE c.AccountID = ?
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.PostID, &comment.CommentID, &comment.Content, &comment.ImageUrl, &comment.Likes, &comment.Dislikes, &comment.AccountID, &comment.AccountUsername, &comment.AccountImageUrl, &comment.CreationDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// DeleteComment function deletes a comment from the database.
func DeleteComment(db *sql.DB, commentID string) error {
	_, err := db.Exec("DELETE FROM comments WHERE commentID = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}

// IncrementNumberOfLikes function increments the number of likes of a comment.
func IncrementNumberOfLikes(db *sql.DB, commentID string) error {
	_, err := db.Exec("UPDATE comments SET likes = likes + 1 WHERE commentID = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}

// IncrementNumberOfDislikes function increments the number of dislikes of a comment.
func IncrementNumberOfDislikes(db *sql.DB, commentID string) error {
	_, err := db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE commentID = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}

// DecrementNumberOfLikes function decrements the number of likes of a comment.
func DecrementNumberOfLikes(db *sql.DB, commentID string) error {
	_, err := db.Exec("UPDATE comments SET likes = likes - 1 WHERE commentID = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}

// DecrementNumberOfDislikes function decrements the number of dislikes of a comment.
func DecrementNumberOfDislikes(db *sql.DB, commentID string) error {
	_, err := db.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE commentID = ?", commentID)
	if err != nil {
		return err
	}
	return nil
}
