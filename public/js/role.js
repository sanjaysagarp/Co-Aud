$(document).ready(function(){
	$('#auditionSubmit').prop('disabled', true);
	
	$('#uploadInput').click(function() {
		$('#auditionFile').trigger('click');
	});

	$('#auditionFile').change(function() {
		var vals = $(this).val();
		val = vals.length ? vals.split('\\').pop() : '';
		$('#uploadInput').val(val);
		if($('#uploadInput').val() =="") {
			$('#auditionSubmit').prop('disabled', true);
		} else {
			$('#auditionSubmit').prop('disabled', false);
		}
	});
	
	$("#formAudition").submit(function(){

		var formData = new FormData($(this)[0]);
		formData.append("id", $('#roleCommentButton').data("id"));
		$.ajax({
			url: "/api/v1/submitAudition/",
			type: 'POST',
			data: formData,
			async: false,
			success: function (data) {
				console.log(data);
				if (data =="rejected") {
					$("#notification").css("display", "block");
					$("#notification").addClass("alert alert-danger");
					$("#notification").html("Filesize too big!");
					$("#notification").fadeOut( 3000 );
					
				}
				if(data == "uploaded") {
					window.location.href = "/role/?id=" + $('#roleCommentButton').data("id")
				}
				
			},
			cache: false,
			contentType: false,
			processData: false
		});
	});

	$.validate({
		form : '#formAudition',
		modules : 'file, toggleDisabled, security',
		disabledFormFilter : 'form.toggle-disabled',
		showErrorDialogs : false,
		onkeyup: true
	});
	
});
