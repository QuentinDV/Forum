// Fonction pour récupérer la valeur d'un cookie par son nom
function getCookieValue(cookieName) {
    var name = cookieName + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var cookieArray = decodedCookie.split(';');
    for (var i = 0; i < cookieArray.length; i++) {
        var cookie = cookieArray[i].trim();
        if (cookie.indexOf(name) === 0) {
            return cookie.substring(name.length, cookie.length);
        }
    }
    return "";
}

// Get the value of the "account" cookie
var accountCookieValue = getCookieValue("account");

// Split the cookie value using the separator "|"
var accountValues = accountCookieValue.split("|");
console.log(accountValues);

// Get the values of the Account fields
var id = accountValues[0];
var email = accountValues[1];
var password = accountValues[2];
var username = accountValues[3];
var imageUrl = accountValues[4];
var isBan = accountValues[5] === 'true';
var isModerator = accountValues[6] === 'true';
var isAdmin = accountValues[7] === 'true';
var creationDate = accountValues[8];

// Insert the values into the HTML elements
document.getElementById("usernameSpan").innerText = username;
document.getElementById("idSpan").innerText = id;
document.getElementById("emailSpan").innerText = email;
document.getElementById("passwordSpan").innerText = password;
// document.getElementById("creationDateSpan").innerText = creationDate;

// Find the img element by its ID and assign the image URL
var profilePicture = document.getElementById("profilePicture");
profilePicture.src = imageUrl;
