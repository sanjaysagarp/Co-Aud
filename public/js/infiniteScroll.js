$(document).ready(function(){

    $(window).scroll(function() {
        if($(window).scrollTop() == $(document).height() - $(window).height()) {
            // ajax call get data from server and append to the div
            $.ajax({
				method: 'POST',
				url: "/api/v1/updateUser/",
				data: {
					displayName: $('[name="displayName"]').val(),
					title: $('[name="title"]').val(),
					personalWebsite: $('[name="personalWebsite"]').val(),
					aboutMe: $('[name="aboutMe"]').val(),
					facebookURL: $('[name="facebookURL"]').val(),
					twitterURL: $('[name="twitterURL"]').val(),
					instagramURL: $('[name="instagramURL"]').val()
				},
				dataType: 'html',
				success: function(data) {
					if(data) {
						console.log(data);
						$("#notification").css("display", "block");
						$("#notification").addClass("alert alert-success");
						$("#notification").html("User profile updated!");
						$("#notification").fadeOut( 3000 );
					} 
					
				},
				failure: function(err) {
					console.log(err);
				}
		});
        }
    });
});