$(document).ready(function(){
    $('#submit-button').affix({
    offset: {
        top: 0,
        bottom: $('body').outerHeight(true) - (50 + $('#navbar-block').outerHeight(true) + $('.content-wrapper').outerHeight(true))
    }
    });
 });