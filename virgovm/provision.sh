#!/bin/bash

apt-get -qq -y update
apt-get -qq -y install vim curl git bzr make g++ gcc


chown vagrant:vagrant /data/gopath
chown vagrant:vagrant /data/gopath/src
chown vagrant:vagrant /data/gopath/src/github.com
chown vagrant:vagrant /data/gopath/src/github.com/racker

mkdir -p /data/O_O
mount --bind /data/gopath/src/github.com/racker/go-agent-endpoint/virgovm/O_O /data/O_O

if [ ! -d "/data/O_O/go" ]; then
  (cd /data/O_O && curl 'https://go.googlecode.com/files/go1.1.linux-amd64.tar.gz' | tar zxv)
fi

cp /data/O_O/conf/profile /home/vagrant/.profile

sudo -u vagrant bash -c 'source /data/O_O/conf/rc; go get -u github.com/songgao/colorgo'
sudo -u vagrant bash -c 'source /data/O_O/conf/rc; go get -u launchpad.net/gocheck'

