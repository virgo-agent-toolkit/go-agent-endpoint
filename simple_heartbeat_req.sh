#!/bin/bash

echo '{"id": 1, "target": "n1", "source": "n2", "params":{"timestamp": 1368231473191}, "method":"heartbeat.post"}' | socat - TCP4:localhost:9876
