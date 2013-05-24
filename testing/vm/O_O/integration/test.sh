#!/bin/bash

/usr/bin/stud -q $(dirname $0)/test.pem -b 127.0.0.1,50050 -f *,50051 --ssl &
stud_pid=$!

endpoint_bin=$1
agent_out=$2

$endpoint_bin ':50050' &
endpoint_pid=$!

echo '
monitoring_id agentA
monitoring_token 0000000000000000000000000000000000000000000000000000000000000000.7777
monitoring_endpoints 127.0.0.1:50051
' > /tmp/cfg
$agent_out/rackspace-monitoring-agent -i --zip $agent_out/virgo-bundle.zip --config /tmp/cfg 1>/dev/null 2>&1 &
agent_pid=$!

sleep 5

kill $stud_pid
kill $endpoint_pid
kill $agent_pid
