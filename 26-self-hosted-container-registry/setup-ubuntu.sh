#!/bin/bash

# install docker engine
sudo apt install docker.io
sudo docker --version

# install docker-compose
# look here for the latest version and arch - https://github.com/docker/compose/releases
sudo curl -L "https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-linux-aarch64" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo docker-compose --version

# create a directory for the registry
cd ~
mkdir -p ./registry/data

# create a password file, install apache2-utils
sudo apt install apache2-utils
htpasswd -Bbn busy bee > ./registry/registry.password

mkdir -p nginx/certs
# put your certs in the nginx/certs directory
