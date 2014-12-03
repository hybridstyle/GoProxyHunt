#!/bin/sh

SPATH=$(cd "$(dirname "$0")"; pwd)
export GOROOT="/macken/go"
SRCPATH=$SPATH/../src
export GOPATH=$SRCPATH


$GOROOT/bin/go build $SRCPATH/GoProxyHunt.go

