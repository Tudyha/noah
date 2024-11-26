#!/bin/sh

cd ../web || exit
npm run build:stage

cd ..

GOOS=linux GOARCH=amd64 go build -o noah cmd/ns/ns.go