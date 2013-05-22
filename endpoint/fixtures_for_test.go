package endpoint

// manually entered for now; TODO: should be generated from .json files by script

const (
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
)
