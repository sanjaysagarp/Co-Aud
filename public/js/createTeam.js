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

$(document).ready(
	function(){
		$("#submitTeam").click(function() {
			var checker = 0;
			if(validateText("teamName")) {
				checker++;
			}
			if(validateText("motto")) {
				checker++;
			}
			$('.teamEmails').each(function () {
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
			if(checker == (2 + (counter))) {
				$("form#teamForm").submit();
			} else {
				return false;
			}
		});
		
		var counter = 1;
		$('#addMoreTeam').on('click', function(e) {
			counter++;
			var newdiv = document.createElement('div');
			newdiv.className = "pair";
			newdiv.innerHTML = "<br><div class='form-group castForm'><input type='text' class='teamEmails form-control' name='teamEmails'></div>";
			document.getElementById("dynamicInput").appendChild(newdiv);
			document.getElementById("removeTeam").style.display = 'inline';
		});
		$('#removeTeam').on('click', function(e) {
			if(counter == 2) {
				counter--;
				var inputs = $(".pair");
				$(inputs[inputs.length - 1]).remove();
				document.getElementById("removeTeam").style.display = 'none';
				return false;
			}
			counter--;
			var inputs = $(".pair");
			 $(inputs[inputs.length - 1]).remove();
		});	
	}
);	