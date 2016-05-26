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
	
	$(document).ajaxSend(function(event, request, settings) {
		$('#loading-indicator').show();
	});

	$(document).ajaxComplete(function(event, request, settings) {
		$('#loading-indicator').hide();
	});
	
	$("#formAudition").submit(function(e){
		e.preventDefault();
		console.log($('#auditionSubmit').data("id"));
		var id = $('#auditionSubmit').data("id");
		var formData = new FormData(this);
		
		formData.append('id', id);
		formData.append('auditionFile', $('input[id="auditionFile"]')[0].files[0]);
		
		console.log(formData);
		$.ajax({
			url: "/api/v1/submitAudition/",
			type: 'POST',
			data: formData,
			contentType: false,
			processData: false,
			cache: false,
			success: function (data) {
				console.log(data);
				if (data =="rejected") {
					$("#notification").css("display", "block");
					$("#notification").addClass("alert alert-danger");
					$("#notification").html("Filesize too big!");
					$("#notification").fadeOut( 3000 );
				}
				if(data == "uploaded") {
					window.location.href = "/auditions/?id=" + id
				}
			}
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
