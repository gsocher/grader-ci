#!/bin/bash
# Script for provisioning a vagrant vm

# install go and setup bashrc
mkdir -p /home/vagrant/bin
curl -sL -o /home/vagrant/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x /home/vagrant/bin/gimme

gimme_command='$(GIMME_GO_VERSION=1.8 /home/vagrant/bin/gimme)'

cat >> /home/vagrant/.bashrc << EOF
eval $gimme_command
export GOPATH=/home/vagrant
alias ci="cd /home/vagrant/src/github.com/dpolansky/grader-ci"
ci
EOF

source /home/vagrant/.bashrc

# install docker https://docs.docker.com/engine/installation/linux/ubuntu/
sudo apt-get update
sudo apt-get install curl \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual \
    apt-transport-https \
    software-properties-common \
    ca-certificates

curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add -

sudo add-apt-repository \
       "deb https://apt.dockerproject.org/repo/ \
       ubuntu-$(lsb_release -cs) \
       main"

sudo apt-get update
sudo apt-get -y install docker-engine

sudo groupadd docker
sudo gpasswd -a vagrant docker
sudo service docker restart
sudo newgrp docker

# install rabbitmq https://www.rabbitmq.com/install-debian.html
echo 'deb http://www.rabbitmq.com/debian/ testing main' |
        sudo tee /etc/apt/sources.list.d/rabbitmq.list

wget -O- https://www.rabbitmq.com/rabbitmq-release-signing-key.asc |
        sudo apt-key add -

sudo apt-get update
sudo apt-get -y install rabbitmq-server

# install sqlite
sudo apt-get install sqlite3 libsqlite3-dev


# build docker images
sh build_docker_images.sh