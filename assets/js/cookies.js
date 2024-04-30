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

// Récupérer la valeur du cookie "account"
var accountCookieValue = getCookieValue("account");

// Diviser la valeur du cookie en utilisant le séparateur "|"
var accountValues = accountCookieValue.split("|");
console.log(accountValues);

// Récupérer les valeurs des champs de l'Account
var id = accountValues[0];
var email = accountValues[1];
var password = accountValues[2];
var username = accountValues[3];
var imageUrl = accountValues[4];
var isBan = accountValues[5];
var isAdmin = accountValues[6];
var creationDate = accountValues[7];

// Insérer les valeurs dans les éléments HTML
document.getElementById("usernameSpan").innerText = username;
document.getElementById("idSpan").innerText = id;
document.getElementById("emailSpan").innerText = email;
document.getElementById("passwordSpan").innerText = password;
document.getElementById("isBanSpan").innerText = isBan;
document.getElementById("isAdminSpan").innerText = isAdmin;
document.getElementById("creationDateSpan").innerText = creationDate;

// Trouver l'élément img par son ID et attribuer l'URL de l'image
var profilePicture = document.getElementById("profilePicture");
profilePicture.src = imageUrl;
