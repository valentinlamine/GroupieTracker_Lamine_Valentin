function PlayPreview() {
    var element = document.getElementById('source');
    if (element.paused) {
        element.play();
    }
    else {
        element.pause();
    }
}