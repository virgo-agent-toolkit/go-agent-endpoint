package endpoint

// manually entered for now; TODO: should be generated from .json files by script

var (
	FIXTURE_heartbeat_post_request = `
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

	FIXTURE_heartbeat_post_response = `
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
	FIXTURE_heartbeat_post_request_invalid_version = `
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

	FIXTURE_handshake_hello_request = `
  {
    "v": "1",
    "id": 0,
    "source": "agentA",
    "target": "endpoint",
    "method": "handshake.hello",
    "params": {
      "token": "MYTOKEN",
      "agent_id": "MYUID",
      "agent_name": "Rackspace Monitoring Agent",
      "process_version": "1.0.0",
      "bundle_version": "1.0.0"
    }
  }
`

	FIXTURE_proactive_test_request = &struct {
		Hello string `json:"hello"`
	}{Hello: "world"}
)
