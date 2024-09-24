#!/bin/sh

cd ../web || exit
npm run build:stage

cd ..

go build -o noah cmd/noah/main.go