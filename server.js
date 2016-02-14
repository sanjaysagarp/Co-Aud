"use strict";

var express = require('express'); //the xpress app framework
var glob = require('glob'); //allows for globbing files by names
var favicon = require('serve-favicon'); //literally oinly serves favicons
var morgan = require('morgan'); //logs things to files 
var bodyParser = require('body-parser'); //parses json, html, etc to html
//var passport = require('passport'); //user authentication

// GOOGLE AND FACEBOOK
var session = require('express-session'); //for user sessions
//var fs = require('fs'); //to read files
var mongoose = require('mongoose');
var config = require('./secret/config'); //configuration file 

//var socket = require('socket.io'); //for back and forth communication
// var db = require('./app/models'); //mongo schemas


var app = express(); //let's get started!

// persistent login sessions
// app.use(passport.initialize());
// app.use(passport.session()); 

var port = process.env.PORT || 80;

// connect to a mongoose db, just add to config.db (ex: config.db + '/posts')
// mongoose.connect(config.db);
// mongoose.connection.on('error', function(err) {
// 	console.error(err);
// });

app.set('views', config.root + 'app/views');
app.set('view engine', 'jade');


//load certificate files
// var pubCert = fs.readFileSync(config.root + 'security/server-cert.pem', 'utf-8');
// var privKey = fs.readFileSync(config.root + 'security/server-pvk.pem', 'utf-8');

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