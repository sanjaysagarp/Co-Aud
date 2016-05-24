$(document).ready(function(){

	$('#role-comment-button').on('click', function(e) {
		e.preventDefault();
		if ($.trim($('#role-comment-content').val()) != "") {
			$.ajax({
					method: 'POST',
					url: "/api/v1/submitRoleComment/",
					data: {
						//disable comment section
						content: $('#role-comment-content').val(),
						id: $('#role-comment-button').data("id")
					},
					dataType: 'html',
					success: function(data) {
						if(data) {
							$('#role-comment-section').prepend(data);
							$('#role-comment-section li:first-child').prepend("<hr>");
							$('#role-comment-content').val("");
							$('#role-comment-content').css("height", "74");
							$('#no-comment-message').remove();
						} 
						
					},
					failure: function(err) {
						console.log(err);
						$("#notification").css("display", "block");
						$("#notification").addClass("alert alert-warning");
						$("#notification").html("Your comment failed to submit, please try again later!");
						$("#notification").fadeOut( 3000 );
						//reenables comment section after failure
						$('#role-comment-button').prop("disabled", false);
						$('#role-comment-content').prop("disabled", false);
					}
			});
		} else {
			$("#notification").css("display", "block");
			$("#notification").addClass("alert alert-warning");
			$("#notification").html("You cannot submit a blank comment");
			$("#notification").fadeOut( 3000 );
		}
	});
	
	// set audition comment link/button to show comment section on click
	$('.auditionCommentButton').on('click', function(e) {
		$('#auditionCommentSection').remove();
		var thisCommentButton = $(this);
		var auditionId = $(thisCommentButton).data("auditionId");
		var userProfilePic = $("#nav-bar-profile-picture").attr('src');
		var html = '<li id="auditionCommentSection" class="comment-submission media"><div class="media-left"><a href="/profile/"><img class="img-profile media-object img-circle" src="' + userProfilePic + '"></a></div><div class="media-body"><form><textarea id="auditionCommentContent" class="form-control" rows="3" placeholder="Leave a comment..."></textarea><button id="submitAuditionCommentButton" type="submit" class="blue-btn-color btn btn-default pull-right" data-id="' + auditionId + '">Post</button></form></div></li>';
		$(html).insertBefore(thisCommentButton.parent());
		
		//hide comment button
		$(thisCommentButton).css("visibility", "hidden");
		
		$('#submitAuditionCommentButton').on('click', function(e) {
			e.preventDefault();
			if ($.trim($('#auditionCommentContent').val()) != "") {
				//disable commenting
				$('#submitAuditionCommentButton').prop("disabled", true);
				$('#auditionCommentContent').prop("disabled", true);
				$.ajax({
						method: 'POST',
						url: "/api/v1/submitAuditionComment/",
						data: {
							content: $('#auditionCommentContent').val(),
							id: $('#submitAuditionCommentButton').data("id")
						},
						dataType: 'html',
						success: function(data) {
							if(data) {
								var auditionCommentSection = $('#submitAuditionCommentButton').parent().parent().parent().parent();
								//display newly added comment
								// auditionCommentSection.append(data);
								$(data).insertBefore(thisCommentButton.parent());
								//TODO: displays comment button again
								$(thisCommentButton).css("visibility", "visible");
								//remove comment section if success
								$('#auditionCommentSection').remove();
							}
						},
						failure: function(err) {
							console.log(err);
							$("#notification").css("display", "block");
							$("#notification").addClass("alert alert-warning");
							$("#notification").html("Your comment failed to submit, please try again later!");
							$("#notification").fadeOut( 3000 );
							//reenables comment section after failure
							$('#submitAuditionCommentButton').prop("disabled", false);
							$('#auditionCommentContent').prop("disabled", false);
						}
				});
			} else {
				$("#notification").css("display", "block");
				$("#notification").addClass("alert alert-warning");
				$("#notification").html("You cannot submit a blank comment");
				$("#notification").fadeOut( 3000 );
			}
		});
	});
});