package database

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type LogInError struct {
	IndentifError bool
	PassWordError bool
}

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
	correctpassword := !checkPassword(password, account.Password)
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

// Fonction pour vérifier si l'identifiant est un email
func IsEmail(identif string) bool {
	return strings.Contains(identif, "@")
}

// Fonction pour récupérer un compte par son email
func GetAccountByEmail(db *sql.DB, email string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT id, email, password, username, ImageUrl, isAdmin, isBan, CreationDate FROM accounts WHERE email = ?", email)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsAdmin, &account.IsBan, &account.CreationDate)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// Fonction pour récupérer un compte par son username
func GetAccountByUsername(db *sql.DB, username string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT id, email, password, username, ImageUrl, isAdmin, isBan, CreationDate FROM accounts WHERE username = ?", username)
	err := row.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsAdmin, &account.IsBan, &account.CreationDate)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func checkPassword(password, hashedPassword string) bool {
	return hashedPassword == hashPasswordSHA256(password)
}
