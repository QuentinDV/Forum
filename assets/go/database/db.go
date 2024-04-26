package database

import "database/sql"

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
		username TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
