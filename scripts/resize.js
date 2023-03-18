function ResizeCover() {
    var cover = document.getElementById("book_cover");
    var description = document.getElementById("book_description");
    cover.style.maxHeight = description.offsetHeight + "px";
}

window.onload = ResizeCover;
window.onresize = ResizeCover;