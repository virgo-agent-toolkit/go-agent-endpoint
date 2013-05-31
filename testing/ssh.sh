#!/bin/bash

cd $(dirname $0)/vm
PORT=$(vagrant ssh-config | grep Port | awk '{print $2}')
USER=$(vagrant ssh-config | grep 'User ' | awk '{print $2}')
KEY=$(vagrant ssh-config | grep IdentityFile | awk '{print $2}' | sed -e 's/^"//'  -e 's/"$//')
HOST=$(vagrant ssh-config | grep HostName | awk '{print $2}')

ssh -q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i "$KEY" -p $PORT $USER@$HOST $@
