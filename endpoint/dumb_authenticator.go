package endpoint

type DumbAuthenticatorDontUseMe byte

func (auth DumbAuthenticatorDontUseMe) Authenticate(agentName string, agentId string, token string, connCtx ConnContext) HandleCode {
	return OK
}
