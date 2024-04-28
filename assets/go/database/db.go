package database

import "database/sql"

type Account struct {
	Id           string
	Email        string
	Password     string
	Username     string
	ImageUrl     string
	IsBan        bool
	IsAdmin      bool
	CreationDate string
}

// Fonction pour créer une nouvelle connexion à la base de données SQLite
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
        isAdmin BOOLEAN NOT NULL DEFAULT 0, 
        CreationDate TEXT NOT NULL
    )`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Fonction pour insérer un nouvel compte dans la base de données
func InsertAccount(db *sql.DB, account Account) error {
	_, err := db.Exec("INSERT INTO accounts (id, email, password, username, ImageUrl, isBan, isAdmin, CreationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		account.Id, account.Email, account.Password, account.Username, account.ImageUrl, account.IsBan, account.IsAdmin, account.CreationDate)
	return err
}

func DeleteAccount(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM accounts WHERE id = ?", id)
	return err
}

func GetAllAccounts(db *sql.DB) ([]Account, error) {
	rows, err := db.Query("SELECT id, email, password, username, ImageUrl, isBan, isAdmin, CreationDate FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err = rows.Scan(&account.Id, &account.Email, &account.Password, &account.Username, &account.ImageUrl, &account.IsBan, &account.IsAdmin, &account.CreationDate)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
