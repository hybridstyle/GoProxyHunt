#!/bin/sh

SPATH=$(cd "$(dirname "$0")"; pwd)
export GOROOT="/usr/local/go"
SRCPATH=$SPATH/../src
export GOPATH=$SPATH/..

cd $SRCPATH
$GOROOT/bin/go build GoProxyHunt.go

