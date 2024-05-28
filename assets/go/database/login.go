package database

// Importing necessary packages
import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type LogInError struct {
	IndentifError bool
	PassWordError bool
}

// RecoverAccount function retrieves an account from the database using the provided identifier and password.
func RecoverAccount(identif, password string) (Account, LogInError, error) {
	var emptyaccount Account
	var account Account
	// Connexion à la base de données
	db, err := ConnectDB("database.db")
	if err != nil {
		return emptyaccount, LogInError{}, err
	}
	defer db.Close()

	// Vérifier si l'identifiant est un email ou un pseudonyme
	if IsEmail(identif) {
		account, err = GetAccountByEmail(db, identif)
	} else {
		account, err = GetAccountByUsername(db, identif)
	}

	if err != nil {
		return emptyaccount, LogInError{}, err
	}

	// Vérifier si le compte existe
	correctpassword := checkPassword(password, account.Password)
	correctidentif := !(account == emptyaccount)

	if !correctpassword && !correctidentif {
		return emptyaccount, LogInError{IndentifError: true, PassWordError: true}, nil

	} else if !correctidentif {
		return emptyaccount, LogInError{IndentifError: true, PassWordError: false}, nil

	} else if !correctpassword {
		return emptyaccount, LogInError{IndentifError: false, PassWordError: true}, nil
	}

	return account, LogInError{}, nil
}

// IsEmail function checks if the provided identifier is an email.
func IsEmail(identif string) bool {
	return strings.Contains(identif, "@")
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
	row := db.QueryRow("SELECT * FROM accounts WHERE username = ?", username)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsModerator, &account.IsAdmin, &account.CreationDate)
	return account, err
}

// checkPassword function checks if the provided password matches the hashed password from the database.
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
