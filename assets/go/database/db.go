package database

// Importing necessary packages
import (
	"database/sql"
	"fmt"
)

// Account struct represents a user account in the system
type Account struct {
	Id           string
	Email        string
	Password     string
	Username     string
	ImageUrl     string
	IsBan        bool
	IsModerator  bool
	IsAdmin      bool
	CreationDate string
}

// ConnectDB function creates a new connection to the SQLite database
func ConnectDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Créer la table accounts si elle n'existe pas déjà
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        username TEXT UNIQUE NOT NULL,
        ImageUrl TEXT NOT NULL,
        isBan BOOLEAN NOT NULL DEFAULT 0,
        isModerator BOOLEAN NOT NULL DEFAULT 0, 
        isAdmin BOOLEAN NOT NULL DEFAULT 0, 
        CreationDate TEXT NOT NULL
    )`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InsertAccount function inserts a new account into the database.
// It takes a database connection and an account as input.
// It returns an error if any.
func InsertAccount(db *sql.DB, account Account) error {
	_, err := db.Exec("INSERT INTO accounts (id, email, password, username, ImageUrl, isBan, isModerator, isAdmin, CreationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		account.Id, account.Email, account.Password, account.Username, account.ImageUrl, account.IsBan, account.IsModerator, account.IsAdmin, account.CreationDate)
	return err
}

// DeleteAccount function deletes an account from the database.
// It takes a database connection and an account ID as input.
// It returns an error if any.
func DeleteAccount(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM accounts WHERE id = ?", id)
	return err
}

// GetAllAccounts function retrieves all accounts from the database.
// It takes a database connection as input.
// It returns a slice of accounts and an error if any.
func GetAllAccounts(db *sql.DB) ([]Account, error) {
	rows, err := db.Query("SELECT id, email, password, username, ImageUrl, isBan, isModerator, isAdmin, CreationDate FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	fmt.Println("test")
	fmt.Println(rows)
	fmt.Println(accounts)
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate); err != nil {
			return nil, err
		}
		fmt.Println(account)
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func ChangeData(db *sql.DB, id string, username string, imageUrl string) error {
	_, err := db.Exec("UPDATE accounts SET username = ?, ImageUrl = ? WHERE id = ?", username, imageUrl, id)
	return err
}

func ChangePassword(db *sql.DB, id string, password string) error {
	_, err := db.Exec("UPDATE accounts SET password = ? WHERE id = ?", password, id)
	return err
}

func BanAccount(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isBan = 1 WHERE id = ?", id)
	return err
}

func UnBanAccount(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isBan = 0 WHERE id = ?", id)
	return err
}

func PromoteToModerator(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isModerator = 1 WHERE id = ?", id)
	return err
}

func DemoteFromModerator(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isModerator = 0 WHERE id = ?", id)
	return err
}

func PromoteToAdmin(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isAdmin = 1 WHERE id = ?", id)
	return err
}

func DemoteFromAdmin(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isAdmin = 0 WHERE id = ?", id)
	return err
}
