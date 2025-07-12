#!/bin/bash

# VARIABLES
INSTALL_DIR = "/usr/bin

# asking super user permission
sudo su

echo "Installing file sharing client"

# install go 
echo "Installing go"
sudo apt install go 

echo "building the package"
go build main.go -o  ${INSTALL_DIR}/fsgoclient

echo "successfully installed the client!!\n for more help to use it type fsgoclient -help"

exit

