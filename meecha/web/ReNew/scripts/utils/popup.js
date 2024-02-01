const popups = document.querySelectorAll(".popup");

popups.forEach(function(element) {
    element.addEventListener("click",function(evt) {
        evt.target.classList.remove('is-show');
    })
});

