$(document).ready(function(){
	var counter = 1;
	$('#submitWork').on('click', function(e) {
		var casts = []; 
		var castRoles = [];
		for(var i = 1; i <= counter; i++) {
			var cast = $('[name="castEmail' + i + '"]').val();
			casts.push(cast);
			var castRole = $('[name="castRole' + i + '"]').val();
			castRoles.push(castRole);
			
		}
		e.preventDefault();
		$.ajax({
				method: 'POST',
				url: "/api/v1/publishWork/",
				data: {
					castList : casts.toString(),
					castRoles : castRoles.toString(),
					title : $('[name="title"]').val(),
					URL : $	('[name="url"]').val(),
					description : $('[name="description"]').val(),
					shortDescription : $('[name="shortDescription"]').val()
				},
				dataType: 'html',
				success: function(data) {	
					if(data) {
						$("#notification").css("display", "block");
						$("#notification").addClass("alert alert-success");
						$("#notification").html("Project has been created!");
						$("#notification").fadeOut( 3000 );
						//location.href = "/seanTest/";
						//window.location.href = "/castings.html/"
						//"/project/id?=" + $('#submitWork').data("id")
					} 
					
				},
				failure: function(err) {
					console.log(err);
				}
		});
	});
	
	$('#addMoreCast').on('click', function(e) {
		counter++;
		var newdiv = document.createElement('div');
		newdiv.innerHTML = "Cast Member " + counter + "<br><input type='text' name='castEmail" + counter + "'> as <input type='text' name='castRole"+ counter +"'>";
		document.getElementById("dynamicInput").appendChild(newdiv);
	
	});
});