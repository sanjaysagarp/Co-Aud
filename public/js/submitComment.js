$(document).ready(function(){

	$('#roleCommentButton').on('click', function(e) {
		e.preventDefault();
		$.ajax({
				method: 'POST',
				url: "/api/v1/submitRoleComment/",
				data: {
					//disable comment section
					content: $('#roleCommentContent').val(),
					id: $('#roleCommentButton').data("id")
				},
				dataType: 'html',
				success: function(data) {
					if(data) {
						//display new comment
						//reenable comment section
						console.log(data);
						$("#notification").css("display", "block");
						$("#notification").addClass("alert alert-success");
						$("#notification").html("Your comment has been submitted!");
						$("#notification").fadeOut( 3000 );
					} 
					
				},
				failure: function(err) {
					console.log(err);
				}
		});
	});
	
	// $('#auditionCommentButton').on('click', function(e) {
	// 	e.preventDefault();
	// 	$.ajax({
	// 			method: 'POST',
	// 			url: "/api/v1/submitAuditionComment/",
	// 			data: {
	// 				//disable comment section
	// 				content: $('[name="auditionCommentContent"]').val(),
	// 				id: $('#auditionCommentButton').data("id")
	// 			},
	// 			dataType: 'html',
	// 			success: function(data) {
	// 				if(data) {
	// 					//display new comment
	// 					//reenable comment section
	// 					console.log(data);
	// 					$("#notification").css("display", "block");
	// 					$("#notification").addClass("alert alert-success");
	// 					$("#notification").html("Your comment has been submitted!");
	// 					$("#notification").fadeOut( 3000 );
	// 				} 
					
	// 			},
	// 			failure: function(err) {
	// 				console.log(err);
	// 			}
	// 	});
	// });
	
	//ON CLICK 
	// <li class="comment-submission media">
	// 	<div class="media-left">
	// 		<a href="/profile/">
	// 			<img class="img-profile media-object img-circle" src="/public/img/default_profile_pic.png">
	// 		</a>
	// 	</div>
	// 	<div class="media-body">
	// 		<form>
	// 			<textarea class="form-control" rows="3" placeholder="Leave a comment..."></textarea>
	// 			<button type="submit" class="btn btn-default pull-right">Post</button>
	// 		</form>
	// 	</div>
	// </li>
	
	
	// $('#auditionCommentButton').on('click', function(e) {
	// 	e.preventDefault();
	// 	$.ajax({
	// 			method: 'POST',
	// 			url: "/api/v1/submitComment/",
	// 			data: {
	// 				content: $('[name="auditionCommentContent"]').val(),
	// 				collection: "roles"
	// 				id: 
	// 			},
	// 			dataType: 'html',
	// 			success: function(data) {
	// 				if(data) {
	// 					console.log(data);
	// 					$("#notification").css("display", "block");
	// 					$("#notification").addClass("alert alert-success");
	// 					$("#notification").html("Your comment has been submitted!");
	// 					$("#notification").fadeOut( 3000 );
	// 				} 
					
	// 			},
	// 			failure: function(err) {
	// 				console.log(err);
	// 			}
	// 	});
	// });
});