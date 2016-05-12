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
		div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback"></span>');
		return true;
	}
}
$(document).ready(
	function(){
		$("#submitWork").click(function() {
			var checker = 0;
			if(validateText("name")) {
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
			if(checker == 6 + counter) {
				$("form#castingForm").submit();
			} else {
				return false;
			}
		});
	}
);

