$(document).ready(function(){
	$('#updateButton').on('click', function(e) {
		e.preventDefault();
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
						window.location.href = "/profile/?id=" + data;
					} 
					
				},
				failure: function(err) {
					console.log(err);
				}
		});
	});
	
	$('input[id="profPic"]').change(function() { 
		var formData = new FormData(this);
		
		formData.append('profPic', $('input[id="profPic"]')[0].files[0]);
		
		//e.preventDefault();
		$.ajax({
			method: 'POST',
			url: "/api/v1/updateUserPicture/",
			data: formData,
			contentType: false,
			processData: false,
			cache: false,
			success: function(data) {
				if(data != "rejected") {
					$('#userImage').attr('src',data);
				} else {
					$("#notification").css("display", "block");
					$("#notification").addClass("alert alert-danger");
					$("#notification").html("Image too big!");
					$("#notification").fadeOut( 3000 );
				}
				
			},
			failure: function(err) {
				console.log(err);
			}
		});
	});
	
	$(document).ajaxSend(function(event, request, settings) {
		$('#loading-indicator').show();
	});

	$(document).ajaxComplete(function(event, request, settings) {
		$('#loading-indicator').hide();
	});
});