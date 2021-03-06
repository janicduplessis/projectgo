Go Chat
=========

[![Build Status](https://travis-ci.org/janicduplessis/projectgo.svg?branch=master)](https://travis-ci.org/janicduplessis/projectgo)

Description
-----------
Go Chat is a chat server written in [go](http://golang.org/) and a web client using [Polymer](http://www.polymer-project.org/).

Prerequisites
----------
* [go](http://golang.org/)
* [git](http://git-scm.com/)
* [mercurial](http://mercurial.selenic.com/)
* [npm](https://www.npmjs.org/)
* [mysql](http://www.mysql.com/)

Installation
----------
1. Get and build the code using

        go get github.com/janicduplessis/projectgo

2. Install build system globally

        npm install bower -g
        npm install grunt -g

3. Navigate to the project root and get other dependencies
		
		cd $GOPATH/src/github.com/janicduplessis/projectgo
        npm install

4. Run using 

        grunt

5. After the first run, a config file will be created: *server.json* edit it with your server information

6. 	Run again!

        grunt 

7. The server should be up and running on localhost:8080 or the address in the config file
