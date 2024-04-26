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