#!/bin/bash

endpoint_bin=$1
agent_out=$2
echo '
monitoring_id agentA
monitoring_token 0000000000000000000000000000000000000000000000000000000000000000.7777
monitoring_endpoints 127.0.0.1:50051
monitoring_upgrade true
' > /tmp/cfg

cmd1="/usr/bin/stud -q $(dirname $0)/test.pem -b 127.0.0.1,50050 -f *,50051 --ssl --write-proxy"
cmd2="$endpoint_bin ':50050' 'localhost:8989'"
cmd3="sudo $agent_out/rackspace-monitoring-agent -i --debug --zip $agent_out/virgo-bundle.zip --config /tmp/cfg"
cmd4="mkdir -p /tmp/upgrading && cd /tmp/upgrading && python2.7 -m SimpleHTTPServer 8989"

pingpong -log "/data/O_O/logs/$(date)" -- "$cmd1" "$cmd2" "$cmd3" "$cmd4"   &
pid=$?

sleep 32
kill $pid
exit 0
