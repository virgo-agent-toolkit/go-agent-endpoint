#!/bin/bash

cd $(dirname $0)/vm
status=$(vagrant status | grep default | sed 's/^default[ \t]*//g')

if [ "$status" == "powered off" ] || [ "$status" == "not created" ]
then
  vagrant up
fi
