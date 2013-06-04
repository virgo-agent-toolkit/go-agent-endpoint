package endpoint

import (
	"sort"
)

type Authenticator interface {
	Authenticate(agentName string, agentId string, token string) bool
}

func authForIterator(auth Authenticator, agentName string, agentId string, token string) HandleCode {
	if auth.Authenticate(agentName, agentId, token) {
		return OK
	} else {
		return FAIL_PASSON
	}
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

func (l *authenticatorList) Iterate(agentName string, agentId string, token string) bool {
	al := *l
	ret := FAIL_PASSON
	for _, item := range al {
		ret = authForIterator(item.authenticator, agentName, agentId, token)
		if !ret.IsPASSON() {
			break
		}
	}
	return ret.IsOK()
}
