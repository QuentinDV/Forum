// script.js
document.addEventListener('DOMContentLoaded', (event) => {
    const createCategoryButton = document.getElementById('createCategoryButton');
    const categoryPopup = document.getElementById('categoryPopup');
    const closeButton = document.getElementsByClassName('close')[0];

    createCategoryButton.onclick = function() {
        categoryPopup.style.display = 'block';
    }

    closeButton.onclick = function() {
        categoryPopup.style.display = 'none';
    }

    window.onclick = function(event) {
        if (event.target == categoryPopup) {
            categoryPopup.style.display = 'none';
        }
    }
});
