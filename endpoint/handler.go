package endpoint

import (
	"sort"
)

type HandleCode int

const (
	OK       = HandleCode(0x0)
	DECLINED = HandleCode(0x1)
	FAIL     = HandleCode(0x2)
)

// General interface for all handlers; should return OK if the request is
// properly handled, DECLINED if the request is not handled yet and should be
// passed on to next handler, or FAIL if there's an error and should stop
// without passing on.
type Handler interface {
	Handle(req *Request, responder *Responder, connContext ConnContext) HandleCode
}

type handlerListItem struct {
	handler Handler

	// the lower the number, the higher the priority, i.e., the earlier it should
	// be executed
	priority int
}

func constructHandlerListItem(handler Handler, priority int) handlerListItem {
	return handlerListItem{handler: handler, priority: priority}
}

type handlerList []handlerListItem

func newHandlerList() *handlerList {
	ret := handlerList(make([]handlerListItem, 0))
	return &ret
}

func (l *handlerList) Len() int { return len(*l) }

func (l *handlerList) Less(i, j int) bool {
	hl := *l
	return hl[i].priority < hl[j].priority // lower number (higher priority) at front
}

func (l *handlerList) Swap(i, j int) {
	hl := *l
	hl[i], hl[j] = hl[j], hl[i]
}

func (hl *handlerList) Push(x handlerListItem) {
	*hl = append(*hl, x)
	sort.Sort(hl)
}

func (l *handlerList) Iterate(req *Request, responder *Responder, connCxt ConnContext) HandleCode {
	hl := *l
	ret := DECLINED
	for _, item := range hl {
		ret = item.handler.Handle(req, responder, connCxt)
		if OK == ret || FAIL == ret {
			break
		}
	}
	return ret
}
