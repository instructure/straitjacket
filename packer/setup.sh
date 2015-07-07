#!/bin/bash
set -e

INSTALL_DIR=/home/ubuntu/straitjacket

mkdir $INSTALL_DIR
tar xvof /home/ubuntu/straitjacket.tar -C $INSTALL_DIR
chown -R ubuntu:ubuntu $INSTALL_DIR

apt-get update
apt-get install -y btrfs-tools
mkfs.btrfs /dev/xvdb
mkdir /var/lib/docker
mount /dev/xvdb /var/lib/docker
echo /dev/xvdb /var/lib/docker btrfs defaults 0 0 >> /etc/fstab
wget -qO- https://get.docker.com/ | sh
usermod -aG docker ubuntu
systemctl enable docker

cd /home/ubuntu/straitjacket

./install_apparmor_profiles

docker build -t cstest languages/csharp
docker build -t jstest languages/nodejs
docker build -t rtest languages/ruby

docker build -t straitjacket .
docker run -d --restart=always -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp -p 8081:8081 straitjacket