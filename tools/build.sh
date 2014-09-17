#!/bin/sh
# Simple script for pull -> build -> run
# Used for github webhook to get latest code and restart the server
# http://fideloper.com/node-github-autodeploy
# Hard coded for my server

go_path="/home/ct/golang"
project_repo="/home/ct/golang/src/github.com/janicduplessis/projectgo"
port_nbr="9898"

export GOPATH=$go_path
cd $project_repo
git pull
pid=$(lsof -t -i:${port_nbr})
if [ ! -z "$pid" ]; then
	kill $pid
fi
go run main.go