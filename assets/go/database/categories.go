package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// Categories struct represents a user account in the system
type Category struct {
	Id            string
	Title         string
	Description   string
	ImageUrl      string
	NomberOfPosts int
	View          int
	Tags          []string
	AccountID     string
	CreationDate  string
}

// ConnectCategoriesDB function creates a new connection to the SQLite database
func ConnectCategoriesDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Creating the categories table if it does not already exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        CategoryId TEXT PRIMARY KEY,
        title TEXT UNIQUE NOT NULL,
        description TEXT NOT NULL,
        ImageUrl TEXT NOT NULL,
		nomberofposts INTEGER DEFAULT 0,
        view INTEGER DEFAULT 0, 
		tags TEXT,
        AccountID TEXT UNIQUE NOT NULL,
        creationDate TEXT NOT NULL
    )`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InsertCategory function inserts a new category into the database.
// It takes a database connection and a category as input.
// It returns an error if any.
// Function to insert a category into the database
func InsertCategory(db *sql.DB, category Category) error {
	// Get the last post ID
	row := db.QueryRow("SELECT CategoryId FROM categories ORDER BY CategoryId DESC LIMIT 1")
	var lastID string
	err := row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If there are no categories in the database, set the last ID to "0"
			lastID = "-1"
		} else {
			// If there's another error, return it
			return err
		}
	}

	// Increment the last ID
	newID := incrementID(lastID)
	fmt.Println(newID)
	category.Id = newID

	tagsStr := strings.Join(category.Tags, ",") // Join tags with a comma
	_, err = db.Exec(`INSERT INTO categories (CategoryId, title, description, ImageUrl, nomberofposts, view, tags, AccountID, creationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		category.Id, category.Title, category.Description, category.ImageUrl, category.NomberOfPosts, category.View, tagsStr, category.AccountID, category.CreationDate)
	return err
}

// DeleteCategory function deletes a category from the database.
// It takes a database connection and a category ID as input.
// It returns an error if any.
func DeleteCategory(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM categories WHERE CategoryId = ?", id)
	return err
}

// GetAllCategories function retrieves all categories from the database.
// It takes a database connection as input.
// It returns a slice of categories and an error if any.
func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		var tagsStr string
		err := rows.Scan(&category.Id, &category.Title, &category.Description, &category.ImageUrl, &category.NomberOfPosts, &category.View, &tagsStr, &category.AccountID, &category.CreationDate)
		if err != nil {
			return nil, err
		}
		category.Tags = strings.Split(tagsStr, ",") // Split tags by comma
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategory function retrieves a category from the database.
func GetCategory(db *sql.DB, id string) (Category, error) {
	row := db.QueryRow("SELECT * FROM categories WHERE CategoryId = ?", id)
	var category Category
	var tagsStr string
	err := row.Scan(&category.Id, &category.Title, &category.Description, &category.ImageUrl, &category.NomberOfPosts, &category.View, &tagsStr, &category.AccountID, &category.CreationDate)
	if err != nil {
		return Category{}, err
	}
	category.Tags = strings.Split(tagsStr, ",") // Split tags by comma
	return category, nil
}

// GetCategoriesByCreator function retrieves all categories from a specific creator.
func GetCategoriesByCreator(db *sql.DB, AccountID string) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories WHERE AccountID = ?", AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		var tagsStr string
		err := rows.Scan(&category.Id, &category.Title, &category.Description, &category.ImageUrl, &category.NomberOfPosts, &category.View, &tagsStr, &category.AccountID, &category.CreationDate)
		if err != nil {
			return nil, err
		}
		category.Tags = strings.Split(tagsStr, ",") // Split tags by comma
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategoriesByTag function retrieves all categories with a specific tag.
func GetCategoriesByTag(db *sql.DB, tag string) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories WHERE tags LIKE ?", "%"+tag+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		var tagsStr string
		err := rows.Scan(&category.Id, &category.Title, &category.Description, &category.ImageUrl, &category.NomberOfPosts, &category.View, &tagsStr, &category.AccountID, &category.CreationDate)
		if err != nil {
			return nil, err
		}
		category.Tags = strings.Split(tagsStr, ",") // Split tags by comma
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategoriesByTitle function retrieves all categories with a specific title.
func GetCategoriesByTitle(db *sql.DB, title string) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories WHERE title LIKE ?", "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		var tagsStr string
		err := rows.Scan(&category.Id, &category.Title, &category.Description, &category.ImageUrl, &category.NomberOfPosts, &category.View, &tagsStr, &category.AccountID, &category.CreationDate)
		if err != nil {
			return nil, err
		}
		category.Tags = strings.Split(tagsStr, ",") // Split tags by comma
		categories = append(categories, category)
	}
	return categories, nil
}

// AddTagToCategory function adds a tag to a category.
func AddTagToCategory(db *sql.DB, categoryID string, tag string) error {
	category, err := GetCategory(db, categoryID)
	if err != nil {
		return err
	}
	category.Tags = append(category.Tags, tag)
	tagsStr := strings.Join(category.Tags, ",") // Join tags with a comma
	_, err = db.Exec("UPDATE categories SET tags = ? WHERE CategoryId = ?", tagsStr, categoryID)
	return err
}

// RemoveTagFromCategory function removes a tag from a category.
func RemoveTagFromCategory(db *sql.DB, categoryID string, tag string) error {
	category, err := GetCategory(db, categoryID)
	if err != nil {
		return err
	}
	for i, t := range category.Tags {
		if t == tag {
			category.Tags = append(category.Tags[:i], category.Tags[i+1:]...)
			break
		}
	}
	tagsStr := strings.Join(category.Tags, ",") // Join tags with a comma
	_, err = db.Exec("UPDATE categories SET tags = ? WHERE CategoryId = ?", tagsStr, categoryID)
	return err
}

// IncrementNumberOfPosts function increments the number of posts in a category.
func IncrementNumberOfPosts(db *sql.DB, categoryID string) error {
	// Get the number of posts in the category
	rows, err := db.Query("SELECT COUNT(*) FROM posts WHERE categoryID = ?", categoryID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return err
		}
	}

	// Update the number of posts in the category
	_, err = db.Exec("UPDATE categories SET nomberofposts = ? WHERE CategoryId = ?", count, categoryID)
	return err
}

// IncrementView function increments the view count of a category.
func IncrementView(db *sql.DB, categoryID string) error {
	_, err := db.Exec("UPDATE categories SET view = view + 1 WHERE CategoryId = ?", categoryID)
	return err
}
