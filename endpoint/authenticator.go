package endpoint

import (
	"errors"
	"sort"
)

var (
	AuthenticationFailed = errors.New("Authentication failed")
)

// General interface for authenticators; should return OK if the agent is authenticated, or FAIL if this authenticator can't authenticate the agen and the auth information should be passed on
type Authenticator interface {
	Authenticate(agentName string, agentId string, token string, connCtx ConnContext) HandleCode
}

type authenticatorListItem struct {
	authenticator Authenticator
	priority      int // the greater the number (priority), the earlier it should be executed
}

func constructAuthenticatorListItem(authenticator Authenticator, priority int) authenticatorListItem {
	return authenticatorListItem{authenticator: authenticator, priority: priority}
}

type authenticatorList []authenticatorListItem

func newAuthenticatorList() *authenticatorList {
	ret := authenticatorList(make([]authenticatorListItem, 0))
	return &ret
}

func (l *authenticatorList) Len() int { return len(*l) }

func (l *authenticatorList) Less(i, j int) bool {
	al := *l
	return al[i].priority > al[j].priority // higher priority at front
}

func (l *authenticatorList) Swap(i, j int) {
	al := *l
	al[i], al[j] = al[j], al[i]
}

func (l *authenticatorList) Push(x authenticatorListItem) {
	*l = append(*l, x)
	sort.Sort(l)
}

func (l *authenticatorList) Iterate(agentName string, agentId string, token string, connCtx ConnContext) HandleCode {
	al := *l
	ret := FAIL
	for _, item := range al {
		ret = item.authenticator.Authenticate(agentName, agentId, token, connCtx)
		if OK == ret {
			break
		}
	}
	return ret
}
