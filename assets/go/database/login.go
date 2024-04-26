package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Fonction pour récupérer un compte par son email
func GetAccountByEmail(db *sql.DB, email string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT id, email, password, username, ImageUrl, isAdmin FROM accounts WHERE email = ?", email)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsAdmin)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// Fonction pour récupérer un compte par son username
func GetAccountByUsername(db *sql.DB, username string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT id, email, password, username, ImageUrl, isAdmin FROM accounts WHERE username = ?", username)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsAdmin)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func checkPassword(password, hashedPassword string) bool {
	return hashedPassword == hashPasswordSHA256(password)
}
