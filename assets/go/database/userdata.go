package database

import (
	"database/sql"
	"fmt"
	"strings"
)

type UserData struct {
	AccountID            string
	SubscribedCategories []string
	LikedPosts           []string
	DisLikedPosts        []string
	SavedPosts           []string
}

// ConnectUserDataDB function creates a new connection to the SQLite database
func ConnectUserDataDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS userdata (
		AccountID TEXT PRIMARY KEY,
		SubscribedCategories TEXT,
		LikedPosts TEXT,
		DisLikedPosts TEXT,
		SavedPosts TEXT
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertUserData(db *sql.DB, userData UserData) error {
	fmt.Println(len(userData.SubscribedCategories))
	SubCategories := strings.Join(userData.SubscribedCategories, ",")
	LikedPosts := strings.Join(userData.LikedPosts, ",")
	DisLikedPosts := strings.Join(userData.DisLikedPosts, ",")
	SavedPosts := strings.Join(userData.SavedPosts, ",")
	_, err := db.Exec("INSERT INTO userdata (AccountID, SubscribedCategories, LikedPosts, DisLikedPosts, SavedPosts) VALUES (?, ?, ?, ?, ?)",
		userData.AccountID, SubCategories, LikedPosts, DisLikedPosts, SavedPosts)
	return err
}

func GetUserData(db *sql.DB, AccountID string) (UserData, error) {
	var userData UserData
	var subCategories, likedPosts, dislikedPosts, savedPosts string
	row := db.QueryRow("SELECT AccountID, SubscribedCategories, LikedPosts, DisLikedPosts, SavedPosts FROM userdata WHERE AccountID = ?", AccountID)
	err := row.Scan(&userData.AccountID, &subCategories, &likedPosts, &dislikedPosts, &savedPosts)
	if err != nil {
		return userData, err
	}

	userData.SubscribedCategories = strings.Split(subCategories, ",")
	userData.LikedPosts = strings.Split(likedPosts, ",")
	userData.DisLikedPosts = strings.Split(dislikedPosts, ",")
	userData.SavedPosts = strings.Split(savedPosts, ",")

	return userData, nil
}

func UpdateUserData(db *sql.DB, userData UserData) error {
	SubCategories := strings.Join(userData.SubscribedCategories, ",")
	LikedPosts := strings.Join(userData.LikedPosts, ",")
	DisLikedPosts := strings.Join(userData.DisLikedPosts, ",")
	SavedPosts := strings.Join(userData.SavedPosts, ",")
	_, err := db.Exec("UPDATE userdata SET SubscribedCategories = ?, LikedPosts = ?, DisLikedPosts = ?, SavedPosts = ? WHERE AccountID = ?",
		SubCategories, LikedPosts, DisLikedPosts, SavedPosts, userData.AccountID)
	return err
}

func DeleteUserData(db *sql.DB, AccountID string) error {
	_, err := db.Exec("DELETE FROM userdata WHERE AccountID = ?", AccountID)
	return err
}

func GetAllUserData(db *sql.DB) ([]UserData, error) {
	rows, err := db.Query("SELECT AccountID, SubscribedCategories, LikedPosts, DisLikedPosts, SavedPosts FROM userdata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userDatas []UserData
	for rows.Next() {
		var userData UserData
		var subCategories, likedPosts, dislikedPosts, savedPosts string
		if err := rows.Scan(&userData.AccountID, &subCategories, &likedPosts, &dislikedPosts, &savedPosts); err != nil {
			return nil, err
		}
		userData.SubscribedCategories = strings.Split(subCategories, ",")
		userData.LikedPosts = strings.Split(likedPosts, ",")
		userData.DisLikedPosts = strings.Split(dislikedPosts, ",")
		userData.SavedPosts = strings.Split(savedPosts, ",")
		userDatas = append(userDatas, userData)
	}
	return userDatas, nil
}

func AddSubscribedCategory(db *sql.DB, AccountID string, categoryID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}

	// Check if categoryID already exists
	for _, id := range userData.SubscribedCategories {
		if id == categoryID {
			return nil // Category already exists, do nothing
		}
	}

	IncrementSubscriber(db, categoryID)

	// Add categoryID to SubscribedCategories
	userData.SubscribedCategories = append(userData.SubscribedCategories, categoryID)
	return UpdateUserData(db, userData)
}

func RemoveSubscribedCategory(db *sql.DB, AccountID string, categoryID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}
	for i, id := range userData.SubscribedCategories {
		if id == categoryID {
			DecrementSubscriber(db, categoryID)
			userData.SubscribedCategories = append(userData.SubscribedCategories[:i], userData.SubscribedCategories[i+1:]...)
			break
		}
	}
	return UpdateUserData(db, userData)
}

func AddLikedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}

	// Check if postID already exists
	for _, id := range userData.LikedPosts {
		if id == postID {
			return nil // Liked post already exists, do nothing
		}
	}

	AddLiketoDB(db, postID)

	// Add postID to LikedPosts
	userData.LikedPosts = append(userData.LikedPosts, postID)
	return UpdateUserData(db, userData)
}

func RemoveLikedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}
	for i, id := range userData.LikedPosts {
		if id == postID {
			RemoveLiketoDB(db, postID)
			userData.LikedPosts = append(userData.LikedPosts[:i], userData.LikedPosts[i+1:]...)
			break
		}
	}
	return UpdateUserData(db, userData)
}

func AddDisLikedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}

	// Check if postID already exists
	for _, id := range userData.DisLikedPosts {
		if id == postID {
			return nil // Disliked post already exists, do nothing
		}
	}

	AddDisliketoDB(db, postID)

	// Add postID to DisLikedPosts
	userData.DisLikedPosts = append(userData.DisLikedPosts, postID)
	return UpdateUserData(db, userData)
}

func RemoveDisLikedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}
	for i, id := range userData.DisLikedPosts {
		if id == postID {
			RemoveDisliketoDB(db, postID)
			userData.DisLikedPosts = append(userData.DisLikedPosts[:i], userData.DisLikedPosts[i+1:]...)
			break
		}
	}
	return UpdateUserData(db, userData)
}

func AddSavedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}

	// Check if postID already exists
	for _, id := range userData.SavedPosts {
		if id == postID {
			return nil // Saved post already exists, do nothing
		}
	}

	// Add postID to SavedPosts
	userData.SavedPosts = append(userData.SavedPosts, postID)
	return UpdateUserData(db, userData)
}

func RemoveSavedPost(db *sql.DB, AccountID string, postID string) error {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return err
	}
	for i, id := range userData.SavedPosts {
		if id == postID {
			userData.SavedPosts = append(userData.SavedPosts[:i], userData.SavedPosts[i+1:]...)
			break
		}
	}
	return UpdateUserData(db, userData)
}

func GetSubscribedCategories(db *sql.DB, AccountID string) ([]string, error) {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return nil, err
	}
	return userData.SubscribedCategories, nil
}

func GetLikedPosts(db *sql.DB, AccountID string) ([]string, error) {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return nil, err
	}
	return userData.LikedPosts, nil
}

func GetDisLikedPosts(db *sql.DB, AccountID string) ([]string, error) {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return nil, err
	}
	return userData.DisLikedPosts, nil
}

func GetSavedPosts(db *sql.DB, AccountID string) ([]string, error) {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return nil, err
	}
	return userData.SavedPosts, nil
}

// IsThisPostLiked function checks if a post is liked by a user
func IsThisPostLiked(db *sql.DB, AccountID string, postID string) bool {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return false
	}
	for _, id := range userData.LikedPosts {
		if id == postID {
			return true
		}
	}
	return false
}

// IsThisPostDisliked function checks if a post is disliked by a user
func IsThisPostDisliked(db *sql.DB, AccountID string, postID string) bool {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return false
	}
	for _, id := range userData.DisLikedPosts {
		if id == postID {
			return true
		}
	}
	return false
}

// IsThisPostSaved function checks if a post is saved by a user
func IsThisPostSaved(db *sql.DB, AccountID string, postID string) bool {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return false
	}
	for _, id := range userData.SavedPosts {
		if id == postID {
			return true
		}
	}
	return false
}

func IsThisCategorySubscribed(db *sql.DB, AccountID string, categoryID string) bool {
	userData, err := GetUserData(db, AccountID)
	if err != nil {
		return false
	}
	for _, id := range userData.SubscribedCategories {
		if id == categoryID {
			return true
		}
	}
	return false
}
