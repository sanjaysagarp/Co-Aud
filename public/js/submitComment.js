$(document).ready(function(){

	$('#roleCommentButton').on('click', function(e) {
		e.preventDefault();
		$.ajax({
				method: 'POST',
				url: "/api/v1/submitRoleComment/",
				data: {
					content: $('[name="roleCommentContent"]').val(),
					id: $('#roleCommentButton').data("id")
				},
				dataType: 'html',
				success: function(data) {
					if(data) {
						//disable comment section
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