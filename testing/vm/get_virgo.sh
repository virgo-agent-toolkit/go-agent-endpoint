#!/bin/bash

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
