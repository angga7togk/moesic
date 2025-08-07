#!/bin/bash

APP_NAME="moesic"
OUTPUT_DIR="dist"

rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR


echo "linux/amd64 building..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-linux
echo "linux/amd64 builded!"

echo "windows/amd64 building..."
GOOS=windows GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-windows.exe
echo "windows/amd64 builded!"

echo "darwin/amd64 building..."
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-macos
echo "darwin/amd64 builded!"

echo "darwin/arm64 building..."
GOOS=darwin GOARCH=arm64 go build -o $OUTPUT_DIR/$APP_NAME-macos-arm64
echo "darwin/arm64 builded!"

echo "done all: '$OUTPUT_DIR'"
