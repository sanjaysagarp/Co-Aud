function validateText(id) {
	if ($("#"+id).val() == "" || $("#"+id).val() == null) {
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
		div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback" aria-hidden="true"></span>');
		div.append('<span id="inputSuccess4Status" class="sr-only">(success)</span>');
		return true;
	}
}

function validateForm() {
	var radios = document.getElementsByName("gender");
	var formValid = false;
	for(var increment = 0; increment < radios.length; increment++) {
		console.log(increment);
		if (radios[increment].checked) {
			var div = $("#genderRadio");    
			div.removeClass("has-error");
			div.addClass("has-success");
			div.addClass("has-feedback");
			$("#glypcn2").remove();
			div.append('<span id="glypcn2" class="glyphicon glyphicon-ok form-control-feedback" aria-hidden="true"></span>');
			div.append('<span id="inputSuccess4Status" class="sr-only">(success)</span>');
			formValid = true;
		}
	}
	if (!formValid) {
		var div = $("#genderRadio");
		div.removeClass("has=success");
		$("#glypcn2").remove();
		div.addClass("has-error has-feedback");
		div.append('<span id="glypcn2" class="glyphicon glyphicon-remove form-control-feedback"></span>');
		return formValid;
	} else {
		return formValid;
	}
}
$(document).ready(function(){
		$("#castSubmission").click(function() {
			var checker = 0;
			if(validateText("name")) {
				checker++;
			}
			if(validateText("date")) {
				checker++;
			}
			if(validateText("age")) {
				checker++;
			}
			if(validateText("traits")) {
				checker++;
			}
			if(validateText("description")) {
				checker++;
			}
			if(validateText("script")) {
				checker++;
			}
			if(validateForm()) {
				checker++;
			}
			if(checker == 7) {
				$("form#roleForm").submit();
			} else {
				return false;
			}
		});
		// $("#roleForm").submit(function(e){
		// 	e.preventDefault();
		// 	var formData = new FormData(this);
		// 	console.log(formData);
		// 	$.ajax({
		// 		url: "/update/audition/?id={{ .Data.role.Id.Hex }}",
		// 		type: 'POST',
		// 		data: formData,
		// 		contentType: false,
		// 		processData: false,
		// 		cache: false,
		// 		success: function (data) {
		// 			console.log(data);
		// 			if (data =="rejected") {
		// 				$("#notification").css("display", "block");
		// 				$("#notification").addClass("alert alert-danger");
		// 				$("#notification").html("Filesize too big!");
		// 				$("#notification").fadeOut( 3000 );
		// 			} else {
		// 				window.location.href = "/auditions/?id=" + data;
		// 			}
		// 		}
		// 	});
		// });
});

