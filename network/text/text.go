package text

import (
	"fmt"
	"strings"
)

type MsgHandler func(string, interface{})

type Processor struct {
	msgInfo map[string]MsgHandler
	defaultMsgHandler MsgHandler
}

func NewProcessor() *Processor {
	processor := &Processor{
		msgInfo: make(map[string]MsgHandler),
	}
	return processor
}

func (t *Processor) Route(msg interface{}, userData interface{}) error {
	raw, ok := msg.(string)
	if !ok {
		return fmt.Errorf("message must be plain text: %v", msg)
	}

	segs := strings.Split(raw, " ")
	handler := t.msgInfo[segs[0]]
	if handler == nil {
		t.defaultMsgHandler(raw, userData)
	} else {
		handler(raw, userData)
	}
	return nil
}

func (t *Processor) Unmarshal(data []byte) (interface{}, error) {
	return string(data), nil
}

func (t *Processor) Marshal(msg interface{}) ([][]byte, error) {
	raw, ok := msg.(string)
	if !ok {
		return nil, fmt.Errorf("message must be plain text: %v", msg)
	}

	return [][]byte{[]byte(raw)}, nil
}

func (t *Processor) SetDefaultHandler(handler MsgHandler) {
	t.defaultMsgHandler = handler
}

func (t *Processor) SetHandler(s string, handler MsgHandler) {
	t.msgInfo[s] = handler
}

