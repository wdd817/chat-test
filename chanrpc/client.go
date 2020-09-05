package chanrpc

import (
	"chat-test/util"
	"errors"
	"fmt"
)

// rpc客户端
type Client struct {
	s                *Server
	pendingAsyncCall int
	chanSyncRet      chan *RetInfo
	ChanAsyncRet     chan *RetInfo
}

// NewClient 创建一个rpc客户端
func NewClient(l int) *Client {
	return &Client{
		ChanAsyncRet: make(chan *RetInfo, l),
	}
}

// Attach 挂载rpc服务器
func (c *Client) Attach(s *Server) {
	c.s = s
}

func (c *Client) Call(id interface{}, arg interface{}) (interface{}, error) {
	f := c.s.functions[id]
	if f == nil {
		return nil, fmt.Errorf("function id %v: function not registered", id)
	}

	err := c.call(&CallInfo{
		f:       f,
		arg:     arg,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return nil, err
	}

	ri := <-c.chanSyncRet
	return ri.ret, ri.err
}

// Async 异步调用
func (c *Client) Async(id interface{}, arg interface{}, cb Callback) {
	if c.pendingAsyncCall >= cap(c.ChanAsyncRet) {
		c.execCb(&RetInfo{
			err: errors.New("too many calls"),
			cb:  cb,
		})
		return
	}

	c.async(id, arg, cb)
	c.pendingAsyncCall++
}

func (c *Client) call(ci *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if block {
		c.s.ChanCall <- ci
	} else {
		select {
		case c.s.ChanCall <- ci:
		default:
			err = errors.New("chanrpc channel full")
		}
	}
	return
}

func (c *Client) async(id interface{}, arg interface{}, cb Callback) {
	f := c.s.functions[id]
	if f == nil {
		c.ChanAsyncRet <- &RetInfo{
			err: fmt.Errorf("function id %v: function not registered", id),
			cb:  cb,
		}
		return
	}

	err := c.call(&CallInfo{
		f:       f,
		arg:     arg,
		cb:      cb,
		chanRet: c.ChanAsyncRet,
	}, false)
	if err != nil {
		c.ChanAsyncRet <- &RetInfo{
			err: err,
			cb:  cb,
		}
		return
	}
}

func (c *Client) execCb(ri *RetInfo) {
	defer util.PrintPanicStack()
	ri.cb(ri.ret, ri.err)
}

// Cb 执行结束后的回调
func (c *Client) Cb(ri *RetInfo) {
	c.pendingAsyncCall--
	c.execCb(ri)
}

// Close 关闭客户端
func (c *Client) Close() {
	for c.pendingAsyncCall > 0 {
		c.Cb(<-c.ChanAsyncRet)
	}
}

// Idle 是否空闲
func (c *Client) Idle() bool {
	return c.pendingAsyncCall == 0
}
