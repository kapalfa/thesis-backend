#!/bin/bash 

#check if mkcert is installed 
if ! command -v mkcert &> /dev/null
then
    echo "mkcert could not be found"
    echo "Please install mkcert and run again"
    exit
fi

#generate key and certificate
echo "Generating key and certificate..."
mkcert -install
mkcert localhost

#start server 
echo "Starting server..."