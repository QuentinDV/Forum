package main

// Importing necessary packages
import (
	"fmt"
	"forum/assets/go/database"
	"forum/assets/go/web"
	"net/http"
	"time"
)

func main() {
	// Creating the database
	database.ConnectUserDB("database.db")
	database.CreateAccount("", "", "Guest", true, true)
	database.CreateAccount("quentin.dassivignon@ynov.com", "Quentin123", "QuentinDV", true, true)
	database.CreateAccount("owandji.dieng@ynov.com", "Owandji123", "OwandjiD", true, true)

	dbcategories, err := database.ConnectCategoriesDB("database.db")
	if err != nil {
		fmt.Println("Error connecting to categories database:", err)
	}
	userdatadb, _ := database.ConnectUserDataDB("database.db")
	database.ConnectCategoriesDB("database.db")
	postdb, _ := database.ConnectPostDB("database.db")
	defer postdb.Close()
	defer dbcategories.Close()

	database.InsertCategory(dbcategories, database.Category{Title: "Video Games", Description: "Video Games in General", ImageUrl: "https://media1.giphy.com/media/3oEjHYlwvUK5p9AIbm/giphy.gif?cid=6c09b9524j6keeuo57nnwcf4q3gtje4k0rkjpy94fujx8vkv&ep=v1_gifs_search&rid=giphy.gif&ct=g", Tags: []string{"videogames"}, AccountID: "QuentinDV", CreationDate: time.Now().Format("2006-01-02 15:04:05")})
	database.InsertCategory(dbcategories, database.Category{Title: "Animal", Description: "Animals in General", ImageUrl: "https://media0.giphy.com/media/uUs14eCA2SBgs/200w.gif", Tags: []string{"videogames"}, AccountID: "QuentinDV", CreationDate: time.Now().Format("2006-01-02 15:04:05")})

	database.CreatePost(postdb, "Fortnite", "Fortnite is a battle royale game", "https://media.tenor.com/gjbq2A5lzDcAAAAM/fortnite-fortnite-dance.gif", "1", "1")
	database.CreatePost(postdb, "Minecraft", "Minecraft is a cube game", "", "1", "1")
	database.CreatePost(postdb, "Dogs", "Dogs are cute", "", "2", "1")

	database.AddSubscribedCategory(userdatadb, "2", "0")
	database.AddSubscribedCategory(userdatadb, "2", "1")
	database.AddLikedPost(userdatadb, "2", "3")
	database.AddLikedPost(userdatadb, "2", "4")
	database.AddLikedPost(userdatadb, "2", "4")

	// Pages
	http.HandleFunc("/", web.Home)
	http.HandleFunc("/home", web.Home)
	http.HandleFunc("/categories", web.Categories)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/admin", web.Admin)
	http.HandleFunc("/userprofile", web.UserProfile)
	http.HandleFunc("/otheruserprofile", web.OtherUserProfile)
	http.HandleFunc("/notfound", web.NotFound)

	// Forms
	http.HandleFunc("/signupform", web.SignUpForm)
	http.HandleFunc("/loginform", web.LoginForm)
	http.HandleFunc("/guestform", web.LogOutForm)
	http.HandleFunc("/logoutform", web.LogOutForm)
	http.HandleFunc("/userprofileform", web.UserProfileForm)

	http.HandleFunc("/banUserform", web.BanForm)
	http.HandleFunc("/deleteUserform", web.DeleteAccountForm)
	http.HandleFunc("/promoteToModeratorform", web.ModeratorForm)
	http.HandleFunc("/promoteToAdminform", web.AdminForm)

	http.HandleFunc("/PfpImageForm", web.PfpWithImageForm)
	http.HandleFunc("/ChangePwForm", web.ChangePwForm)

	// Elements
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("./assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("./assets/img"))))

	// Links
	fmt.Println("\nPlay : http://localhost:8080/admin")
	fmt.Println("\nPlay : http://localhost:8080/categories")
	fmt.Println("\nPlay : http://localhost:8080/home")
	fmt.Println("\nPlay : http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
