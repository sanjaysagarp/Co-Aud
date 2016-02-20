"use strict";

var express = require('express'); //the express app framework
var glob = require('glob'); //allows for globbing files by names
var favicon = require('serve-favicon'); //serves favicons
var morgan = require('morgan'); //logs things to files 
var bodyParser = require('body-parser'); //parses json, html, etc to html
var session = require('express-session'); //for user sessions
var mongoose = require('mongoose'); //mongodb connections
var config = require('./secret/config'); //configuration file 



var app = express();

var port = process.env.PORT || 80;

app.set('views', config.root + 'app/views');
app.set('view engine', 'jade');
app.use(express.static(__dirname  + '/public'));

var env = process.env.NODE_ENV || 'development';
app.locals.ENV = env;
app.locals.ENV_DEVELOPMENT = env == 'development';

app.use(favicon('public/img/favicon.ico'));
app.use(morgan('dev'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: true}));
app.use(express.static(config.root + '/public'));

//Grabs all the controllers in the folder, and adds them to the controllers
//list. Synchronous function.
var controllers = glob.sync(config.root + '/app/controllers/*.js');
controllers.forEach(function assignController(controller) {
	require(controller)(app);
});

// lets get started!
app.listen(port);
console.log('server is listening...');