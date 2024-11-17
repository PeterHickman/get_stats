#!/bin/sh

BINARY='/usr/local/bin'
APP=get_stats

echo "Building $APP"
go build -ldflags="-s -w" $APP.go

echo "Installing $APP to $BINARY"
install $APP $BINARY

echo "Removing the build"
rm $APP