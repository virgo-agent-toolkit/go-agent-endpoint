#!/bin/bash

endpoint_bin=$1
agent_out=$2
echo '
monitoring_id agentA
monitoring_token 0000000000000000000000000000000000000000000000000000000000000000.7777
monitoring_endpoints 127.0.0.1:50051
' > /tmp/cfg

cmd1="/usr/bin/stud -q $(dirname $0)/test.pem -b 127.0.0.1,50050 -f *,50051 --ssl --write-proxy"
cmd2="$endpoint_bin ':50050'"
cmd3="$agent_out/rackspace-monitoring-agent -i --zip $agent_out/virgo-bundle.zip --config /tmp/cfg"

pingpong -log "/data/O_O/logs/$(date)" -- "$cmd1" "$cmd2" "$cmd3"   &
pid=$?

sleep 5
kill $pid
exit 0
