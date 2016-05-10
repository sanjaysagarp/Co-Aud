function validateText(id) {
	if ($("#"+id).val() == "") {
		var div = $("#"+id).closest("div");
		div.removeClass("has=success");
		$("#glypcn"+id).remove();
		div.addClass("has-error has-feedback");
		div.append('<span id="glypcn'+id + '" class="glyphicon glyphicon-remove form-control-feedback"></span>');
		return false;
	} else {
		var div = $("#"+id).closest("div");
		div.removeClass("has-error");
		div.addClass("has-success");
		div.addClass("has-feedback");
		$("#glypcn"+id).remove();
		div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback"></span>');
		return true;
	}
}

function validateYoutube(id) {
	var url = $("#"+id).val();
	var isyouTubeUrl = /^(https?:\/\/)?([^\/]*\.)?youtube\.com\/watch\?([^]*&)?v=\w+(&[^]*)?/i.test(url);
	if(!isyouTubeUrl) {
		var div = $("#"+id).closest("div");
		div.removeClass("has=success");
		$("#glypcn"+id).remove();
		div.addClass("has-error has-feedback");
		div.append('<span id="glypcn'+id + '" class="glyphicon glyphicon-remove form-control-feedback"></span>');
		return false;	
	} else {
		var div = $("#"+id).closest("div");
		div.removeClass("has-error");
		div.addClass("has-success");
		div.addClass("has-feedback");
		$("#glypcn"+id).remove();
		div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback"></span>');
		return true;
	}
}

$(document).ready(
	function(){
		$("#submitWork").click(function() {
			var checker = 0;
			if(validateText("title")) {
				checker++;
			}
			if(validateYoutube("url")) {
				checker++;
			}
			if(validateText("shortDescription")) {
				checker++;
			}
			if(validateText("description")) {
				checker++;
			}
			console.log($("#castEmail[]").val());
			if(checker == 4) {
				$("form#workForm").submit();
			} else {
				return false;
			}
		});
		
	var counter = 1;
	// $('#submitWork').on('click', function(e) {
	// 	var casts = []; 
	// 	var castRoles = [];
	// 	for(var i = 1; i <= counter; i++) {
	// 		var cast = $('[name="castEmail' + i + '"]').val();
	// 		casts.push(cast);
	// 		var castRole = $('[name="castRole' + i + '"]').val();
	// 		castRoles.push(castRole);
			
	// 	}
	// 	e.preventDefault();
	// 	$.ajax({
	// 			method: 'POST',
	// 			url: "/api/v1/publishWork/",
	// 			data: {
	// 				castList : casts.toString(),
	// 				castRoles : castRoles.toString(),
	// 				title : $('[name="title"]').val(),
	// 				URL : $	('[name="url"]').val(),
	// 				description : $('[name="description"]').val(),
	// 				shortDescription : $('[name="shortDescription"]').val()
	// 			},
	// 			dataType: 'html',
	// 			success: function(data) {	
	// 				if(data) {
	// 					$("#notification").css("display", "block");
	// 					$("#notification").addClass("alert alert-success");
	// 					$("#notification").html("Project has been created!");
	// 					$("#notification").fadeOut( 3000 );
	// 					//location.href = "/seanTest/";
	// 					//window.location.href = "/castings.html/"
	// 					//"/project/id?=" + $('#submitWork').data("id")
	// 				} 
					
	// 			},
	// 			failure: function(err) {
	// 				console.log(err);
	// 			}
	// 	});
	// });
	
	$('#addMoreCast').on('click', function(e) {
		counter++;
		var newdiv = document.createElement('div');
		newdiv.innerHTML = "Cast Member " + counter + "<br><input type='text' name='castEmail[]'> as <input type='text' name='castRole[]'>";
		document.getElementById("dynamicInput").appendChild(newdiv);
	
	});
});