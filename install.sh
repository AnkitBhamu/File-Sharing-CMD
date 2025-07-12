#!/bin/bash

# Exit on error
set -e

# VARIABLES
INSTALL_DIR="/usr/bin"
GO_VERSION="1.24.4"
CURRENT_DIR=$(pwd)

echo "Installing file sharing client..."

# Ensure script is run with sudo/root privileges
if [[ $EUID -ne 0 ]]; then
   echo "Please run this script with sudo: sudo $0"
   exit 1
fi

# Install Go
echo "Installing Go..."
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C ${CURRENT_DIR} -xzf go${GO_VERSION}.linux-amd64.tar.gz
export PATH=$PATH:${CURRENT_DIR}/go/bin


# Build the package
echo "Building the package..."
go build  -o ${INSTALL_DIR}/fsgoclient  main.go

rm go${GO_VERSION}.linux-amd64.tar.gz
rm -rf ${CURRENT_DIR}/go
echo "Successfully installed the client! To get help, type: fsgoclient -help"
