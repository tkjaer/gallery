document.addEventListener('keydown', function(event) {
    const currentHash = window.location.hash;
    const match = currentHash.match(/^#img-(\d+)$/);

    if (match) {
        const currentIndex = parseInt(match[1], 10);

        if (event.key === 'ArrowRight') {
            const nextIndex = currentIndex - 1;
            window.location.hash = nextIndex > 0 ? `#img-${nextIndex}` : '#';
        } else if (event.key === 'ArrowLeft') {
            const nextIndex = currentIndex + 1;
            /* stop at the last image */
            if (document.getElementById(`img-${nextIndex}`)) {
                window.location.hash = `#img-${nextIndex}`;
            } else {
                window.location.hash = '';
            }
           øwindow.location.hash = `#img-${nextIndex}`;
        } else if (event.key === 'Escape') {
            window.location.hash = '';
        }
    }
});


let touchStartX = 0;
let touchEndX = 0;

document.addEventListener('touchstart', function(event) {
    touchStartX = event.changedTouches[0].screenX;
});

document.addEventListener('touchend', function(event) {
    touchEndX = event.changedTouches[0].screenX;
    const currentHash = window.location.hash;
    const match = currentHash.match(/^#img-(\d+)$/);

    if (match) {
        const currentIndex = parseInt(match[1], 10);

        if (touchStartX - touchEndX > 50) { // Swipe left
            const nextIndex = currentIndex + 1;
            window.location.hash = `#img-${nextIndex}`;
        } else if (touchEndX - touchStartX > 50) { // Swipe right
            const nextIndex = currentIndex - 1;
            window.location.hash = nextIndex > 0 ? `#img-${nextIndex}` : '#';
        }
    }
});