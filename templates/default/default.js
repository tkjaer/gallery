document.addEventListener('keydown', function(event) {
    const currentHash = window.location.hash;
    const match = currentHash.match(/^#(.+)$/);
    const lightboxImages = document.getElementById('lightbox_images');
    const lightboxHrefs = lightboxImages.querySelectorAll('a');
    const lightboxIndex = Array.from(lightboxHrefs).map(element => element.id);

    if (match) {
        imageName = match[1];
        const currentIndex = lightboxIndex.indexOf(imageName);

        if (event.key === 'ArrowRight') {
            const nextIndex = currentIndex + 1;
            if (lightboxIndex[nextIndex]) {
                window.location.hash = `${lightboxIndex[nextIndex]}`;
            } else {
                window.location.hash = '';
            }
        } else if (event.key === 'ArrowLeft') {
            const nextIndex = currentIndex - 1;
            if (lightboxIndex[nextIndex]) {
                window.location.hash = `${lightboxIndex[nextIndex]}`;
            } else {
                window.location.hash = '';
            }
        } else if (event.key === 'Escape') {
            window.location.hash = '';
        }   
    }
});


let touchstartx = 0;
let touchendx = 0;

document.addeventlistener('touchstart', function(event) {
    touchstartx = event.changedtouches[0].screenx;
});

document.addeventlistener('touchend', function(event) {
    touchendx = event.changedtouches[0].screenx;
    const currenthash = window.location.hash;
    const match = currenthash.match(/^#(.+)$/);
    const lightboximages = document.getelementbyid('lightbox_images');
    const lightboxhrefs = lightboximages.queryselectorall('a');
    const lightboxindex = array.from(lightboxhrefs).map(element => element.id);

    if (match) {
        imagename = match[1];
        const currentindex = lightboxindex.indexof(imagename);

        if (touchstartx - touchendx > 50) { // swipe left
            const nextindex = currentindex + 1;
            if (lightboxindex[nextindex]) {
                window.location.hash = `${lightboxindex[nextindex]}`;
            } else {
                window.location.hash = '';
            }
        } else if (touchendx - touchstartx > 50) { // swipe right
            const nextindex = currentindex - 1;
            if (lightboxindex[nextindex]) {
                window.location.hash = `${lightboxindex[nextindex]}`;
            } else {
                window.location.hash = '';
            }
        }
    }
});