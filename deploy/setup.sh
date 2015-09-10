#!/bin/bash
set -e

INSTALL_DIR=/home/ubuntu/straitjacket

# write a blank config file, to be replaced by cloud-init
touch /etc/straitjacket.env

rm -rf $INSTALL_DIR > /dev/null
mkdir $INSTALL_DIR
tar xvof /home/ubuntu/straitjacket.tar -C $INSTALL_DIR
chown -R ubuntu:ubuntu $INSTALL_DIR

apt-get update
apt-get install -y btrfs-tools
if ! mount | grep -q /var/lib/docker; then
  mkfs.btrfs /dev/xvdb
  mkdir /var/lib/docker
  mount /dev/xvdb /var/lib/docker
  echo /dev/xvdb /var/lib/docker btrfs defaults 0 0 >> /etc/fstab
  wget -qO- https://get.docker.com/ | sh
  usermod -aG docker ubuntu
  systemctl enable docker
fi

cd $INSTALL_DIR

./straitjacket-setup.sh

docker build -t straitjacket .
docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp --rm straitjacket --test
