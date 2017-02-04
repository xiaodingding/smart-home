(function($) {
    "use strict"; // Start of use strict

    var owner = 'e154';

    $(document).ready(function () {
        $('#smart-home-server').release({repo: 'smart-home'});
        $('#smart-home-configurator').release({repo: 'smart-home-configurator'});
        $('#smart-home-node').release({repo: 'smart-home-node'})
    });

    //===============================================================
    // LATEST RELEASE
    //===============================================================

    $.fn.release = function(opt) {
        var href = 'https://github.com/'+owner+'/'+opt.repo+'/releases';
        //TODO complite idea
        // $.getJSON("https://api.github.com/repos/"+owner+"/"+opt.repo+"/releases/latest").done(function (release) {
        //     href = release.assets[0].browser_download_url;
        // });
        $(this).attr('href', href);
    };

    //===============================================================
    // ETC
    //===============================================================

    // Disable empty links in docs examples
    $('[href="#"]').click(function (e) {
        e.preventDefault()
    });

    // jQuery for page scrolling feature - requires jQuery Easing plugin
    $('a.page-scroll').bind('click', function(event) {
        var $anchor = $(this);
        $('html, body').stop().animate({
            scrollTop: ($($anchor.attr('href')).offset().top - 50)
        }, 1250, 'easeInOutExpo');
        event.preventDefault();
    });

    $('body').scrollspy({
        target: '#docsNavbarContent',
        offset: 150
    });

    // Closes the Responsive Menu on Menu Item Click
    $('.navbar-collapse ul li a').click(function() {
        $('.navbar-toggle:visible').click();
    });

    // Offset for Main Navigation
    $('#mainNav').affix({
        offset: {
            top: 100
        }
    });

    $('#docsNavbarContent').affix({
        offset: {
            top: 200
        }
    });

    // Initialize and Configure Scroll Reveal Animation
    window.sr = ScrollReveal();
    sr.reveal('.sr-icons', {
        duration: 600,
        scale: 0.3,
        distance: '0px'
    }, 200);
    sr.reveal('.sr-button', {
        duration: 1000,
        delay: 200
    });
    sr.reveal('.sr-contact', {
        duration: 600,
        scale: 0.3,
        distance: '0px'
    }, 300);

    // Initialize and Configure Magnific Popup Lightbox Plugin
    $('.popup-gallery').magnificPopup({
        delegate: 'a',
        type: 'image',
        tLoading: 'Loading image #%curr%...',
        mainClass: 'mfp-img-mobile',
        gallery: {
            enabled: true,
            navigateByImgClick: true,
            preload: [0, 1] // Will preload 0 - before current, and 1 after the current image
        },
        image: {
            tError: '<a href="%url%">The image #%curr%</a> could not be loaded.'
        }
    });

})(jQuery); // End of use strict

(function () {
    'use strict';

    anchors.options.placement = 'left';
    anchors.add('.docs-section > h1, .docs-section > h2, .docs-section > h3, .docs-section > h4, .docs-section > h5')
}());
