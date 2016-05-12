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
		div.removeClass("has-success");
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
			$('.castEmails').each(function () {
				if (this.value == null || this.value == "") {
					var div = $(this).closest("div");
					div.removeClass("has=success");
					//$("#glypcn"+id).remove();
					div.addClass("has-error has-feedback");
					//div.append('<span id="glypcn'+id + '" class="glyphicon glyphicon-remove form-control-feedback"></span>');
				} else {
					var div = $(this).closest("div");
					div.removeClass("has-error");
					div.addClass("has-success");
					div.addClass("has-feedback");
					checker++;
					//$("#glypcn"+id).remove();
					//div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback"></span>');
				}
			})
			$('.castRoles').each(function () {
				if (this.value == null || this.value == "") {
					var div = $(this).closest("div");
					div.removeClass("has=success");
					//$("#glypcn"+id).remove();
					div.addClass("has-error has-feedback");
					//div.append('<span id="glypcn'+id + '" class="glyphicon glyphicon-remove form-control-feedback"></span>');
				} else {
					var div = $(this).closest("div");
					div.removeClass("has-error");
					div.addClass("has-success");
					div.addClass("has-feedback");
					checker++;
					//$("#glypcn"+id).remove();
					//div.append('<span id="glypcn'+id +'" class="glyphicon glyphicon-ok form-control-feedback"></span>');
				}
			})
			if(checker == (4 + (counter * 2))) {
				$("form#workForm").submit();
			} else {
				return false;
			}
		});
		
		var counter = 1;
		$('#addMoreCast').on('click', function(e) {
			counter++;
			var newdiv = document.createElement('div');
			newdiv.className = "pair";
			newdiv.innerHTML = "Cast Member " + counter + "<br><div class='form-group castForm'><input type='text' class='castEmails form-control' name='castEmail[]'></div> as <div class='form-group castForm'><input type='text' name='castRole[]' class='castRoles form-control'></div>";
			document.getElementById("dynamicInput").appendChild(newdiv);
			document.getElementById("removeCast").style.display = 'inline';
		});
		$('#removeCast').on('click', function(e) {
			if(counter == 2) {
				counter--;
				var inputs = $(".pair");
				$(inputs[inputs.length - 1]).remove();
					document.getElementById("removeCast").style.display = 'none';
				return false;
			}
			counter--;
			var inputs = $(".pair");
			 $(inputs[inputs.length - 1]).remove();
		});	
	}
);	