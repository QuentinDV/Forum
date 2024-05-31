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

window.id = accountValues[0][1];
window.email = accountValues[1];
window.password = accountValues[2];
window.username = accountValues[3];
window.imageUrl = accountValues[4];
window.isBan = accountValues[5] === 'true';
window.isModerator = accountValues[6] === 'true';
window.isAdmin = accountValues[7] === 'true';
window.creationDate = accountValues[8];

// Find the img element by its ID and assign the image URL
var profilePicture = document.getElementById("profilePicture");
profilePicture.src = imageUrl;


// Insert the values into the HTML elements
document.getElementById("usernameSpan").innerText = username;
document.getElementById("emailSpan").innerText = email;
document.getElementById("creationDateSpan").innerText = creationDate;


