#!/bin/sh

mkdir -p ./bin
go build -o ./bin/pshunt ./cmd/pshunt/main.go
