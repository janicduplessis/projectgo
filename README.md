projectgo
=========

Description
-----------
CT is a chat server written in go and a web client using Polymer.

Prerequisites
----------
go
git
mercurial
npm
mysql

Installation
----------
1. Get and build the code using 
	$ go get github.com/janicduplessis/projectgo
	$ go install github.com/janicduplessis/projectgo

2. Install build system globally
	$ npm install bower -g
	$ npm install grunt -g

3. Get other dependencies
	$ npm install
	$ bower install

4. Run using $ grunt and then close it

5. After the first run, a config file will be created: server.json edit it with your server information

6. $ grunt again!

7. The server should be up and running on localhost:8080 or the address in the config file
