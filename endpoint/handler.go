package endpoint

import (
	"encoding/json"
	"sort"
)

type HandleCode int

const (
	OK   = HandleCode(0x0)
	FAIL = HandleCode(0x1)
)

// General interface to all handlers; should return OK if the request is properly handled, or FAIL if the request is not handled yet and should be passed on.
type Handler interface {
	Handle(*request, *json.Encoder, ConnContext) HandleCode
}

type handlerListItem struct {
	handler  Handler
	priority int // the lower the number, the higher the priority, i.e., the earlier it should be executed
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

func (l *handlerList) Iterate(req *request, enc *json.Encoder, connCxt ConnContext) HandleCode {
	hl := *l
	ret := FAIL
	for _, item := range hl {
		ret = item.handler.Handle(req, enc, connCxt)
		if OK == ret {
			break
		}
	}
	return ret
}
