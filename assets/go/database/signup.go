package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	Id       string
	Email    string
	Password string
	Username string
}

// Fonction pour insérer un nouvel compte dans la base de données
func InsertAccount(db *sql.DB, account Account) error {
	_, err := db.Exec("INSERT INTO accounts (id, email, password, username) VALUES (?, ?, ?, ?)",
		account.Id, account.Email, account.Password, account.Username)
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

func CreateAccount(email, password, username string) (string, error) {
	// Connexion à la base de données
	db, err := ConnectDB("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Vérifier si l'email est déjà pris
	emailTaken, err := IsEmailTaken(db, email)
	if err != nil {
		return "", err
	}
	if emailTaken {
		return "email", nil
	}

	// Vérifier si le pseudonyme est déjà pris
	usernameTaken, err := IsUsernameTaken(db, username)
	if err != nil {
		return "", err
	}
	if usernameTaken {
		return "username", nil
	}

	// Exemple d'utilisation : Création et insertion d'un nouveau compte
	newAccount := Account{
		Id:       "1",
		Email:    email,
		Password: hashPasswordSHA256(password),
		Username: username,
	}
	err = InsertAccount(db, newAccount)
	if err != nil {
		log.Fatal(err)
	}

	return "", nil
}

// Function for hash a password
func hashPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

// func checkPassword(password, hashedPassword string) bool {
// 	return hashedPassword == hashPasswordSHA256(password)
// }
