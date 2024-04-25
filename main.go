package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Page Principale du Forum
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/home.html")
}

func main() {
	// Pages
	http.HandleFunc("/home", Home)

	// Elements
	http.Handle("/assets/static/", http.StripPrefix("/assets/static/", http.FileServer(http.Dir("./assets/static"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Liens
	fmt.Println("\nPlay : http://localhost:8080/home")
	http.ListenAndServe(":8080", nil)
}

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
