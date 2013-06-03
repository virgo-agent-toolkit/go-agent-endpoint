package endpoint

import ()

type controller struct {
}

func (c *controller) Authenticate(agentName string, agentId string, token string) bool {
	return true
}
