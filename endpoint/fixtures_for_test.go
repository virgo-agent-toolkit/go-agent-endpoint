package endpoint

// manually entered for now; TODO: should be generated from .json files by script

var (
	FIXTUREHeartbeatPostRequest = `
{
  "v": "1",
  "id": 1,
  "source": "agentA",
  "target": "endpoint",
  "method": "heartbeat.post",
  "params": {
    "timestamp": 1325645515246
  }
}
`

	FIXTUREHeartbeatPostResponse = `
{
  "v": "1",
  "id": 1,
  "source": "endpoint",
  "target": "agentA",
  "result": {
    "timestamp": 1325645515246
  }
}
`
	FIXTUREHeartbeatPostRequestInvalidVersion = `
{
    "v": "2147483647",
    "id": 1,
    "source": "endpoint",
    "target": "agentA",
    "method": "heartbeat.post",
    "params": {
        "timestamp": 1325645515246
    }
}
`

	FIXTUREHandshakeHelloRequest = `
  {
    "v": "1",
    "id": 0,
    "source": "agentA",
    "target": "endpoint",
    "method": "handshake.hello",
    "params": {
      "token": "MYTOKEN",
      "agentId": "MYUID",
      "agentName": "Rackspace Monitoring Agent",
      "processVersion": "1.0.0",
      "bundleVersion": "1.0.0"
    }
  }
`

	FIXTUREProactiveTestRequest = &struct {
		Hello string `json:"hello"`
	}{Hello: "world"}
)
