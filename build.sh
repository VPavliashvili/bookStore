#!/bin/sh

swag init -d cmd/api/,api/resource/system/,api/resource/books/
go build -C ./cmd/api/ -v -o ../../main
