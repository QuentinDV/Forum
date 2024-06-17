package database

// Importing necessary packages
import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
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

// ConnectUserDB function creates a new connection to the SQLite database
func ConnectUserDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Creating the accounts table if it does not already exist
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
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func ChangeData(db *sql.DB, id string, username string, imageUrl string) error {
	_, err := db.Exec("UPDATE accounts SET username = ?, ImageUrl = ? WHERE id = ?", username, imageUrl, id)
	return err
}

func ChangePassword(db *sql.DB, id string, username string, oldPassword string, newPassword string) error {
	var password string
	err := db.QueryRow("SELECT password FROM accounts WHERE id = ?", id).Scan(&password)
	if err != nil {
		return err
	}

	if !checkPassword(password, oldPassword) {
		return errors.New("invalid old password")
	}

	newPassword, err = hashPasswordBcrypt(newPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE accounts SET password = ? WHERE id = ?", newPassword, id)
	return err
}

func ChangePWFORCED(db *sql.DB, id string, newPassword string) error {
	newPassword, err := hashPasswordBcrypt(newPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE accounts SET password = ? WHERE id = ?", newPassword, id)
	return err

}

func ShowPassword(db *sql.DB, id string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM accounts WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

// ChangeImageUrl updates the image URL of a user account in the database.
func ChangeImageUrl(db *sql.DB, id string, imageUrl string) error {
	_, err := db.Exec("UPDATE accounts SET ImageUrl = ? WHERE id = ?", imageUrl, id)
	return err
}

// BanAccount function bans an account in the database.
func BanAccount(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isBan = 1 WHERE id = ?", id)
	return err
}

// UnBanAccount function unbans an account in the database.
func UnBanAccount(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isBan = 0 WHERE id = ?", id)
	return err
}

// PromoteToModerator function promotes an account to moderator in the database.
func PromoteToModerator(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isModerator = 1 WHERE id = ?", id)
	return err
}

// DemoteFromModerator function demotes an account from moderator in the database.
func DemoteFromModerator(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isModerator = 0 WHERE id = ?", id)
	return err
}

// PromoteToAdmin function promotes an account to admin in the database.
func PromoteToAdmin(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isAdmin = 1 WHERE id = ?", id)
	return err
}

// DemoteFromAdmin function demotes an account from admin in the database.
func DemoteFromAdmin(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE accounts SET isAdmin = 0 WHERE id = ?", id)
	return err
}

// SaveFile function saves a file to the disk.
func SaveFile(filename string, data io.Reader) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	return nil
}

// CountFiles function retrieves the number of files in a directory and returns it.
func CountFiles(directory string) int {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}
	return len(files)
}

// GetAccount function retrieves an account from the database by ID.
// It takes a database connection and an account ID as input.
// It returns an account and an error if any.
func GetAccountbyID(db *sql.DB, id string) (Account, error) {
	var account Account
	err := db.QueryRow("SELECT id, email, password, username, ImageUrl, isBan, isModerator, isAdmin, CreationDate FROM accounts WHERE id = ?", id).Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// GetAccountByEmail function retrieves an account from the database using the provided email.
func GetAccountByEmail(db *sql.DB, email string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT * FROM accounts WHERE email = ?", email)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate)
	return account, err
}

// GetAccountByUsername function retrieves an account from the database using the provided username.
func GetAccountByUsername(db *sql.DB, username string) (Account, error) {
	var account Account
	fmt.Println("username:", username)
	row := db.QueryRow("SELECT * FROM accounts WHERE username = ?", username)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate)
	return account, err
}

// GetAccountByUsername function retrieves an account from the database by username.
// It takes a database connection and a username as input.
// It returns an account and an error if any.
func GetUserProfileByUsername(db *sql.DB, username string) (Account, error) {
	var account Account
	err := db.QueryRow("SELECT id, email, password, username, ImageUrl, isBan, isModerator, isAdmin, CreationDate FROM accounts WHERE username = ?", username).Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// New function to copy the default profile picture
func CopyDefaultProfilePicture(newID string) (string, error) {
	sourcePath := "./assets/img/pfp/default.png"
	destinationPath := "./assets/img/pfp/" + newID + ".png"

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", err
	}

	return destinationPath, nil
}