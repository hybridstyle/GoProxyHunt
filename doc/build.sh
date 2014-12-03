#!/bin/sh

export GOROOT="/usr/local/go"
GITPATH=/macken/GoProxyHunt
export GOPATH=$GITPATH

cd $GITPATH
git pull
cd src

$GOROOT/bin/go build GoProxyHunt.go

#TODO
