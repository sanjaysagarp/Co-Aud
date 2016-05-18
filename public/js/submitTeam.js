// $(document).ready(function(){
// 	$('#submitTeam').on('click', function(e) {
// 		e.preventDefault();
// 		$.ajax({
// 				method: 'POST',
// 				url: "/api/v1/submitTeam/",
// 				data: {
// 					teamName: $('[name="teamName"]').val(),
// 					motto: $('#role-comment-content').val(),
// 					id: $('#submitTeam').data("id")
// 				},
// 				dataType: 'html',
// 				success: function(data) {
// 					if(data) {
// 						window.location.href = "/profile/edit"
// 					} 	
// 				},
// 				failure: function(err) {
// 					console.log(err);
// 				}
// 		});
		
// 	});
		
// });