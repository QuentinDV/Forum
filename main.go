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
	database.CreateAccount("", "", "Guest", false, false)
	database.CreateAccount("quentin.dassivignon@ynov.com", "Quentin123", "QuentinDV", true, true)
	database.CreateAccount("owandji.dieng@ynov.com", "Owandji123", "OwandjiD", true, true)

	database.ConnectCategoriesDB("db/database.db")
	database.ConnectUserDataDB("db/database.db")
	database.ConnectPostDB("db/database.db")
	database.ConnectCommentsDB("db/database.db")

	// Définition des gestionnaires de routes HTTP
	setupHTTPHandlers()

	// Démarrage du serveur HTTPS
	startHTTPSServer()
}

func setupHTTPHandlers() {
	// Pages
	http.HandleFunc("/", web.LogIn)
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/admin", web.Admin)
	http.HandleFunc("/notfound", web.NotFound)

	http.HandleFunc("/user/", web.UserProfileHandler)
	http.HandleFunc("/category/", web.CategoryPageHandler)
	http.HandleFunc("/post/", web.PostPageHandler)
	http.HandleFunc("/createcategory", web.CreateCategory)

	// Forms
	http.HandleFunc("/signupform", web.SignUpForm)
	http.HandleFunc("/loginform", web.LoginForm)
	http.HandleFunc("/guestform", web.GuestForm)
	http.HandleFunc("/logoutform", web.LogOutForm)

	http.HandleFunc("/createpostform", web.CreatePostForm)
	http.HandleFunc("/createcategoryform", web.CreateCategoryForm)
	http.HandleFunc("/createcommentform", web.CreateCommentForm)

	http.HandleFunc("/postlikeform", web.LikeForm)
	http.HandleFunc("/postdislikeform", web.DislikeForm)
	http.HandleFunc("/commentlikeform", web.LikeCommentForm)
	http.HandleFunc("/commentdislikeform", web.DislikeCommentForm)

	http.HandleFunc("/addviewform", web.AddViewForm)
	http.HandleFunc("/savedpostform", web.SavePostForm)
	http.HandleFunc("/subscribecategoryform", web.SubscribeCategoryForm)
	http.HandleFunc("/sortinghomeform", web.SortingHomePostsForm)
	http.HandleFunc("/resethomesortingform", web.ResetHomeSortingForm)

	http.HandleFunc("/banUserform", web.BanForm)
	http.HandleFunc("/promoteToModeratorform", web.ModeratorForm)
	http.HandleFunc("/promoteToAdminform", web.AdminForm)
	http.HandleFunc("/deleteUserform", web.DeleteAccountForm)
	// http.HandleFunc("/deletecategoryform", web.DeleteCategoryForm)
	http.HandleFunc("/deletecommentform", web.DeleteCommentForm)
	http.HandleFunc("/deletepostform", web.DeletePostForm)

	http.HandleFunc("/PfpImageForm", web.PfpWithImageForm)
	http.HandleFunc("/ChangePwForm", web.ChangePwForm)
	http.HandleFunc("/resetpfpform", web.ResetPfpForm)

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
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}
