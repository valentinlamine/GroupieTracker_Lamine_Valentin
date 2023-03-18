function ResizeCover() {
    var cover = document.getElementById("book_cover");
    var description = document.getElementById("book_description");
    //Si la taille d'écran est supérieur à 991px
    if (window.innerWidth > 991) {
        cover.style.maxHeight = description.offsetHeight + "px";
    } else {
        cover.style.maxHeight = "300px";
    }
}

window.onload = ResizeCover;
window.onresize = ResizeCover;