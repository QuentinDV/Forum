package data

import (
	"database/sql"
	"fmt"
)

func db() {
	// Ouvrir la connexion à la base de données
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de la base de données:", err)
		return
	}
	defer db.Close() // Defer la fermeture de la connexion à la base de données

	// Création de la table
	createTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT,
            email TEXT
        )
    `
	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Erreur lors de la création de la table:", err)
		return
	}
	fmt.Println("Table créée avec succès.")
}
