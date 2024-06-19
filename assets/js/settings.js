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

document.addEventListener('DOMContentLoaded', function() {
    const deleteButtons = document.querySelectorAll('.deleteButton');
    const deleteSound = document.getElementById('deleteSound');

    deleteButtons.forEach(button => {
        button.addEventListener('click', function(event) {
            deleteSound.play();
        });
    });
});

// Get the modal
var modal = document.getElementById("imageModal");

// Get the image and insert it inside the modal - use its "alt" text as a caption
var modalImg = document.getElementById("imgModal");
var captionText = document.getElementById("caption");

document.querySelectorAll('.post img').forEach(img => {
    img.onclick = function(){
        modal.style.display = "block";
        modalImg.src = this.src;
        captionText.innerHTML = this.alt;
    }
});

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("closeModal")[0];

// When the user clicks on <span> (x), close the modal
span.onclick = function() { 
    modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}


document.addEventListener('DOMContentLoaded', () => {
    const deleteButtons = document.querySelectorAll('button[id^="deleteButton"]');
    const deleteSound = document.getElementById('deleteSound');

    deleteButtons.forEach(button => {
        button.addEventListener('click', () => {
            console.log('Button clicked');
            deleteSound.play().catch(error => {
                console.error('Error playing sound:', error);
            });
        });
    });
});