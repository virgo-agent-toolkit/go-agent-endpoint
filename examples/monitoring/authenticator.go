package main

import (
	"github.com/racker/go-agent-endpoint/endpoint"
)

type helloResult struct {
	HeartbeatInterval string `json:"heartbeat_interval"`
	EntityID          string `json:"entity_id"`
	Channel           string `json:"channel"`
}

type authenticator byte

// Authenticate authenticates any agent and returns a heartbeat interval, a
// fake entity ID and "stable" channel. It's used along with endpoint.Hub
func (auth authenticator) Authenticate(agentName string, agentID string, token string, responder *endpoint.Responder, connCtx endpoint.ConnContext) endpoint.HandleCode {
	var result helloResult
	result.HeartbeatInterval = "1000"
	result.EntityID = "fake-entity-id-asdfghjkl"
	result.Channel = "stable"
	responder.Respond(result, nil)

	// return OK for any agent
	return endpoint.OK
}
