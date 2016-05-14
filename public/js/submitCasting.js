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

// function validateForm() {
// 	var radios = document.getElementsByName("gender");
// 	var formValid = false;
// 	var i = 0;
// 	while (!formValid && i < radios.length) {
// 		if (radios[i].checked) {
// 			formValid = true;
// 			i++;  	 
// 			var div = $("#genderRadio");    
// 			div.removeClass("has-error");
// 			div.addClass("has-success");
// 			div.addClass("has-feedback");
// 			$("#glypcn"+id).remove();
// 			div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback" aria-hidden="true"></span>');
// 			div.append('<span id="inputSuccess4Status" class="sr-only">(success)</span>');
// 		}
// 	}
// 	if (!formValid) {
// 		var div = $("#genderRadio");
// 		div.removeClass("has=success");
// 		$("#glypcn"+id).remove();
// 		div.addClass("has-error has-feedback");
// 		div.append('<span id="glypcn2" class="glyphicon glyphicon-remove form-control-feedback"></span>');
// 		return formValid;
// 	} 
// }
$(document).ready(
	function(){
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
			// if(validateForm()) {
			// 	checker++;
			// }
			if(checker == 7) {
				$("form#castingForm").submit();
			} else {
				return false;
			}
		});
	}
);

