#!/bin/sh

swag init -d cmd/api/,api/resource/system/,api/resource/books/
go build -C ./cmd/api/ -v -o ../../main -ldflags "-X main.compileDate=`date +%Y/%m/%d:%H:%M.%S`"
