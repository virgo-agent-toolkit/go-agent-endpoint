#!/bin/bash

apt-get -qq -y update
apt-get -qq -y install vim curl git bzr make g++ gcc


chown vagrant:vagrant /data/gopath
chown vagrant:vagrant /data/gopath/src
chown vagrant:vagrant /data/gopath/src/github.com
chown vagrant:vagrant /data/gopath/src/github.com/racker

mkdir -p /data/O_O
mount --bind /data/gopath/src/github.com/racker/go-agent-endpoint/testing/vm/O_O /data/O_O

if [ ! -d "/data/O_O/go" ]; then
  (cd /data/O_O && curl -s 'https://go.googlecode.com/files/go1.1.linux-amd64.tar.gz' | tar zxv)
fi

cp /data/O_O/conf/profile /home/vagrant/.profile

sudo -u vagrant bash -c 'source /data/O_O/conf/rc; go get -u github.com/songgao/colorgo'
sudo -u vagrant bash -c 'source /data/O_O/conf/rc; go get -u launchpad.net/gocheck'


# Clone virgo repo to /data/virgo. The reasons of doing this in provisioning instead of just using recursively cloned submodule are:
#     1. It seems that current configure in virgo/base doesn't handle well when virgo itself is a submodule
#     2. The only thing here that matters is the SHA1. This could prevent hacking and commiting from submodule
rm -rf /data/virgo
cd /data/gopath/src/github.com/racker/go-agent-endpoint
sha1=$(git submodule status | grep 'testing/virgo' | awk '{print $1}' | sed 's/^[-+]//g')
git clone git://github.com/racker/virgo /data/virgo
cd /data/virgo
git checkout $sha1
git submodule update --init --recursive
chown -R vagrant:vagrant /data/virgo
