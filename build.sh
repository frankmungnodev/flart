#!/bin/bash

VERSION=$(grep "Version =" internal/version/version.go | cut -d '"' -f 2)
BINARY_NAME="flart"
BUILD_DIR="dist"

# Create build directory
mkdir -p $BUILD_DIR

# Build for different platforms
GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/${BINARY_NAME}_${VERSION}_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/${BINARY_NAME}_${VERSION}_darwin_arm64
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/${BINARY_NAME}_${VERSION}_linux_amd64
GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/${BINARY_NAME}_${VERSION}_windows_amd64.exe

# Create archives
cd $BUILD_DIR
tar -czf ${BINARY_NAME}_${VERSION}_darwin_amd64.tar.gz ${BINARY_NAME}_${VERSION}_darwin_amd64
tar -czf ${BINARY_NAME}_${VERSION}_darwin_arm64.tar.gz ${BINARY_NAME}_${VERSION}_darwin_arm64
tar -czf ${BINARY_NAME}_${VERSION}_linux_amd64.tar.gz ${BINARY_NAME}_${VERSION}_linux_amd64
zip ${BINARY_NAME}_${VERSION}_windows_amd64.zip ${BINARY_NAME}_${VERSION}_windows_amd64.exe