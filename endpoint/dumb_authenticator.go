package endpoint

// DumbHelloResult is the Result for handshake.hello response
type DumbHelloResult struct {
	HeartbeatInterval string `json:"heartbeat_interval"`
}

// DumbAuthenticatorDontUseMe implements Authenticator interface; is an
// authenticator that accepts any agent
type DumbAuthenticatorDontUseMe byte

// Authenticate Simply returns OK
func (auth DumbAuthenticatorDontUseMe) Authenticate(agentName string, agentID string, token string, responder *Responder, connCtx ConnContext) HandleCode {
	responder.Respond(DumbHelloResult{HeartbeatInterval: "1000"}, nil)
	return OK
}
