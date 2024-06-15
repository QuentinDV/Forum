package main

// Importing necessary packages
import (
	"fmt"
	"forum/assets/go/database"
	"forum/assets/go/web"
	"net/http"
)

func main() {
	// Creating the database
	database.ConnectUserDB("db/database.db")
	database.CreateAccount("", "", "Guest", true, true)
	database.CreateAccount("quentin.dassivignon@ynov.com", "Quentin123", "QuentinDV", true, true)
	database.CreateAccount("owandji.dieng@ynov.com", "Owandji123", "OwandjiD", true, true)

	database.ConnectCategoriesDB("db/database.db")
	userdatadb, _ := database.ConnectUserDataDB("db/database.db")
	dbcategories, err := database.ConnectCategoriesDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
	database.ConnectPostDB("db/database.db")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}

	err = database.InsertCategory(dbcategories, database.Category{Title: "Video Games", Description: "Video Games in General", ImageUrl: "https://media1.giphy.com/media/3oEjHYlwvUK5p9AIbm/giphy.gif?cid=6c09b9524j6keeuo57nnwcf4q3gtje4k0rkjpy94fujx8vkv&ep=v1_gifs_search&rid=giphy.gif&ct=g", Tags: []string{"videogames"}, AccountID: "QuentinDV"})
	if err != nil {
		fmt.Println("Error inserting category:", err)
	}
	err = database.InsertCategory(dbcategories, database.Category{Title: "Animal", Description: "Animals in General", ImageUrl: "https://media0.giphy.com/media/uUs14eCA2SBgs/200w.gif", Tags: []string{"animal"}, AccountID: "QuentinDV"})
	if err != nil {
		fmt.Println("Error inserting category:", err)
	}

	// err = database.CreatePost(postdb, "Fortnite", "Fortnite is a battle royale game", "https://media.tenor.com/gjbq2A5lzDcAAAAM/fortnite-fortnite-dance.gif", "0", "1")
	// if err != nil {
	// 	fmt.Println("Error inserting Post:", err)
	// }

	// err = database.CreatePost(postdb, "Minecraft", "Minecraft is a cube game", "", "0", "1")
	// if err != nil {
	// 	fmt.Println("Error inserting Post:", err)
	// }
	// err = database.CreatePost(postdb, "Dogs", "Dogs are cute", "", "1", "1")
	// if err != nil {
	// 	fmt.Println("Error inserting Post:", err)
	// }

	database.AddSubscribedCategory(userdatadb, "2", "0")
	database.AddSubscribedCategory(userdatadb, "2", "1")
	// database.RemoveSubscribedCategory(userdatadb, "2", "0")

	database.AddLikedPost(userdatadb, "2", "1")
	database.AddLikedPost(userdatadb, "2", "2")
	database.AddLikedPost(userdatadb, "2", "3")
	database.AddLikedPost(userdatadb, "2", "4")
	database.AddLikedPost(userdatadb, "2", "5")
	database.AddLikedPost(userdatadb, "2", "6")

	database.AddDisLikedPost(userdatadb, "2", "2")
	database.AddDisLikedPost(userdatadb, "2", "1")
	database.AddDisLikedPost(userdatadb, "2", "4")
	database.AddDisLikedPost(userdatadb, "2", "3")
	database.AddDisLikedPost(userdatadb, "2", "5")
	database.AddDisLikedPost(userdatadb, "2", "6")

	// Définition des gestionnaires de routes HTTP
	setupHTTPHandlers()

	// Démarrage du serveur HTTPS
	startHTTPSServer()
}

func setupHTTPHandlers() {
	// Pages
	http.HandleFunc("/", web.LogIn)
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/admin", web.Admin)
	http.HandleFunc("/myprofile", web.UserProfile)
	http.HandleFunc("/notfound", web.NotFound)

	http.HandleFunc("/user/", web.UserProfileHandler)
	http.HandleFunc("/category/", web.CategoryPageHandler)
	http.HandleFunc("/post/", web.PostPageHandler)
	http.HandleFunc("/createcategory", web.CreateCategory)

	// Forms
	http.HandleFunc("/signupform", web.SignUpForm)
	http.HandleFunc("/loginform", web.LoginForm)
	http.HandleFunc("/guestform", web.LogOutForm)
	http.HandleFunc("/logoutform", web.LogOutForm)

	http.HandleFunc("/userprofileform", web.UserProfileForm)
	http.HandleFunc("/createcategoryform", web.CreateCategoryForm)
	http.HandleFunc("/likeform", web.LikeForm)
	http.HandleFunc("/dislikeform", web.DislikeForm)
	http.HandleFunc("/addviewform", web.AddViewForm)

	http.HandleFunc("/banUserform", web.BanForm)
	http.HandleFunc("/deleteUserform", web.DeleteAccountForm)
	http.HandleFunc("/promoteToModeratorform", web.ModeratorForm)
	http.HandleFunc("/promoteToAdminform", web.AdminForm)

	http.HandleFunc("/PfpImageForm", web.PfpWithImageForm)
	http.HandleFunc("/ChangePwForm", web.ChangePwForm)

	// Elements statiques (CSS, JS, images)
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("./assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))
}

func startHTTPSServer() {
	// server := &http.Server{
	// 	Addr: ":443",
	// 	TLSConfig: &tls.Config{
	// 		MinVersion: tls.VersionTLS12,
	// 	},
	// }

	// fmt.Println("\nServeur démarré à : https://localhost:443")
	// fmt.Println("\nServeur démarré à : https://localhost:443/home")

	// err := server.ListenAndServeTLS("cert.pem", "key.pem")
	// if err != nil {
	// 	log.Fatalf("Erreur lors du démarrage du serveur HTTPS: %v", err)
	// }
	// Links
	fmt.Println("\nPlay : http://localhost:8080/admin")
	fmt.Println("\nPlay : http://localhost:8080/categories")
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}
