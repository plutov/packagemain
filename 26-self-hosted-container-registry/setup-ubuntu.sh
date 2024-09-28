#!/bin/bash

sudo apt update

# install docker engine and docker-compose
sudo snap install docker
docker --version
docker-compose --version

# create a directory for the registry
cd ~
mkdir -p ./registry/data

# install htpasswd
sudo apt install apache2-utils
# create a password file
htpasswd -Bbn busy bee > ./registry/registry.password

# put your certificate and key in the certs directory
mkdir -p ./nginx/certs
