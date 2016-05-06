$(document).ready(function(){
    $('#submit-button').affix({
    offset: {
        top: 0,
        bottom: -$('footer').outerHeight(true)
    }
    });
 });