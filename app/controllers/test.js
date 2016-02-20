//test router
var express = require('express');
var router = express.Router();

module.exports = function(app) {
	app.use('/', router);
};

router.get('/test', function(req, res, next) {
	res.render('template/test', {
		title : "im temmie"
	});
});
