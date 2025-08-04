#!/bin/bash

APP_NAME="moesic"
SRC_FILE="main.go"
OUTPUT_DIR="dist"

# Buat folder output
mkdir -p $OUTPUT_DIR

# Build untuk Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-linux $SRC_FILE

# Build untuk Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-windows.exe $SRC_FILE

# Build untuk macOS (64-bit Intel)
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-macos $SRC_FILE

echo "✅ Build selesai! File ada di folder '$OUTPUT_DIR'"
