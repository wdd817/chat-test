package chanrpc

import (
	"errors"
	"fmt"
	"runtime"

	"chat-test/conf"
	"chat-test/log"
)

type Client struct {
	s       *Server
	pending int
	ChanRet chan *RetInfo
}

func NewClient(l int) *Client {
	return &Client{
		ChanRet: make(chan *RetInfo, l),
	}
}

func (c *Client) Attach(s *Server) {
	c.s = s
}

func (c *Client) Async(id interface{}, arg interface{}, cb Callback) {
	if c.pending >= cap(c.ChanRet) {
		c.execCb(&RetInfo{
			err: errors.New("too many calls"),
			cb:  cb,
		})
		return
	}

	c.async(id, arg, cb)
	c.pending++
}

func (c *Client) call(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	select {
	case c.s.ChanCall <- ci:
	default:
		err = errors.New("chanrpc channel full")
	}
	return
}

func (c *Client) async(id interface{}, arg interface{}, cb Callback) {
	f := c.s.functions[id]
	if f == nil {
		c.ChanRet <- &RetInfo{
			err: fmt.Errorf("function id %v: function not registered", id),
			cb:  cb,
		}
		return
	}

	err := c.call(&CallInfo{
		arg:     arg,
		cb:      cb,
		chanRet: c.ChanRet,
	})
	if err != nil {
		c.ChanRet <- &RetInfo{
			err: err,
			cb:  cb,
		}
		return
	}
}

func (c *Client) execCb(ri *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	ri.cb(ri.ret, ri.err)
}

func (c *Client) Cb(ri *RetInfo) {
	c.pending--
	c.execCb(ri)
}

func (c *Client) Close() {
	for c.pending > 0 {
		c.Cb(<-c.ChanRet)
	}
}

func (c *Client) Idle() bool {
	return c.pending == 0
}
