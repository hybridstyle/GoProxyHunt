#!/bin/sh

export GOROOT="/usr/local/go"
GITPATH=/macken/GoProxyHunt
export GOPATH=$GITPATH

cd $GITPATH
$GOROOT/bin/go build GoProxyHunt.go

#TODO
