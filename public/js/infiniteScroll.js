$(document).ready(function() {
	var win = $(window);
	var doc = $(document);
    var page = 1;

	// Each time the user scrolls
	win.scroll(function() {
		// Vertical end reached?
		if (doc.height() - win.height() == win.scrollTop()) {
			// New row
			var roles = $('#home-page-roles');
            $.ajax({
				method: 'POST',
				url: "/api/v1/getRole/",
				data: {
					page: page
				},
				dataType: 'html',
				success: function(data) {
					if(data) {
						roles.append(data);
					}
                    page++;
				},
				failure: function(err) {
					console.log(err);
				}
		    });
		}
	});
});