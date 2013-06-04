package endpoint

import (
	"encoding/json"
	"sort"
)

type HandleCode int

const (
	OK          = HandleCode(0x01)
	OK_PASSON   = HandleCode(0x00)
	FAIL        = HandleCode(0x11)
	FAIL_PASSON = HandleCode(0x10)
)

func (code HandleCode) IsOK() bool {
	return 0x0 == code&0x10
}

func (code HandleCode) IsPASSON() bool {
	return 0x0 == code&0x1
}

type Handler interface {
	Handle(*request, *json.Encoder, *json.Decoder) HandleCode
}

type handlerListItem struct {
	handler  Handler
	priority int // the greater the number (priority), the earlier it should be executed
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
	return hl[i].priority > hl[j].priority // higher priority at front
}

func (l *handlerList) Swap(i, j int) {
	hl := *l
	hl[i], hl[j] = hl[j], hl[i]
}

func (hl *handlerList) Push(x handlerListItem) {
	*hl = append(*hl, x)
	sort.Sort(hl)
}

func (l *handlerList) Iterate(req *request, enc *json.Encoder, dec *json.Decoder) HandleCode {
	hl := *l
	ret := FAIL_PASSON
	for _, item := range hl {
		ret = item.handler.Handle(req, enc, dec)
		if !ret.IsPASSON() {
			break
		}
	}
	return ret
}
