package endpoint

type dumbAuthenticator byte

func (auth dumbAuthenticator) Authenticate(agentName string, agentId string, token string) bool {
	return true
}
