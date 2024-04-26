package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Fonction pour insérer un nouvel compte dans la base de données
func InsertAccount(db *sql.DB, account Account) error {
	_, err := db.Exec("INSERT INTO accounts (id, email, password, username, ImageUrl, isAdmin) VALUES (?, ?, ?, ?, ?, ?)",
		account.Id, account.Email, account.Password, account.Username, account.ImageUrl, account.IsAdmin)
	return err
}

// Fonction pour vérifier si un email est déjà pris dans la base de données
func IsEmailTaken(db *sql.DB, email string) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE email = ?", email)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Fonction pour vérifier si un pseudonyme est déjà pris dans la base de données
func IsUsernameTaken(db *sql.DB, username string) (bool, error) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE username = ?", username)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateAccount(email, password, username, imageUrl string, isAdmin bool) ([]string, error) {
	// Connexion à la base de données
	db, err := ConnectDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the last ID from the database
	var lastID string
	row := db.QueryRow("SELECT id FROM accounts ORDER BY id DESC LIMIT 1")
	err = row.Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			lastID = "0"
		} else {
			return []string{""}, err
		}
	}

	// Increment the last ID
	newID := incrementID(lastID)

	// Vérifier si l'email est déjà pris
	emailTaken, err := IsEmailTaken(db, email)
	if err != nil {
		return []string{""}, err
	}

	// Vérifier si le pseudonyme est déjà pris
	usernameTaken, err := IsUsernameTaken(db, username)
	if err != nil {
		return []string{""}, err
	}

	if emailTaken && usernameTaken {
		return []string{"email", "username"}, nil
	}
	if emailTaken {
		return []string{"email"}, nil
	} else if usernameTaken {
		return []string{"username"}, nil
	}

	// Exemple d'utilisation : Création et insertion d'un nouveau compte
	newAccount := Account{
		Id:       newID,
		Email:    email,
		Password: hashPasswordSHA256(password),
		Username: username,
		ImageUrl: imageUrl,
		IsAdmin:  isAdmin,
	}
	err = InsertAccount(db, newAccount)
	if err != nil {
		log.Fatal(err)
	}

	return []string{""}, nil
}

// Function to increment the ID
func incrementID(lastID string) string {
	var id int
	_, err := fmt.Sscanf(lastID, "%d", &id)
	if err != nil {
		log.Fatal(err)
	}
	id++
	return fmt.Sprintf("%d", id)
}

// Function for hash a password
func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

