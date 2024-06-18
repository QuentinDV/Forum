// Sélectionner le bouton d'ouverture du menu
const openMenuButton = document.getElementById('openMenu');

// Sélectionner le conteneur du menu
const menu = document.getElementById('menu');

// Ajouter un écouteur d'événement pour le clic sur le bouton
openMenuButton.addEventListener('click', function() {
    // Basculer la classe pour afficher ou masquer le menu
    if (menu.style.display === 'block') {
        menu.style.display = 'none';
    } else {
        menu.style.display = 'block';
    }
});

function toggleForm() {
    var form = document.getElementById("createpost");
    if (form.style.display === "none" || form.style.display === "") {
        form.style.display = "block";
    } else {
        form.style.display = "none";
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const likeButton = document.getElementById('likeButton');
    const likeCount = document.getElementById('likeCount');
    let liked = false;

    likeButton.addEventListener('click', () => {
        liked = !liked;
        likeButton.classList.toggle('liked', liked);
        likeCount.textContent = parseInt(likeCount.textContent) + (liked ? 1 : -1);

        // Update the form action
        const likeForm = document.getElementById('likeForm');
        likeForm.submit();
    });
});

function toggleFormcate() {
    var form = document.getElementById("createcate");
    if (form.style.display === "none" || form.style.display === "") {
        form.style.display = "block";
    } else {
        form.style.display = "none";
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const createcate = document.getElementById('createcate');
    const closePopup = document.querySelector('.close');

    closePopup.addEventListener('click', toggleFormcate);

    window.addEventListener('click', (event) => {
        if (event.target === createcate) {
            toggleFormcate();
        }
    });
});
