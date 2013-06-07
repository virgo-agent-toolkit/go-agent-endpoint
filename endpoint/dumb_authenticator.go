package endpoint

// DumbAuthenticatorDontUseMe implements Authenticator interface; is an
// authenticator that accepts any agent
type DumbAuthenticatorDontUseMe byte

// Authenticate Simply returns OK
func (auth DumbAuthenticatorDontUseMe) Authenticate(agentName string, agentID string, token string, connCtx ConnContext) HandleCode {
	return OK
}
