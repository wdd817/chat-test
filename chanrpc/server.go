package chanrpc

import (
	"chat-test/log"
	"errors"
	"fmt"
)

type Server struct {
	// 每类消息对应的处理函数
	functions map[interface{}]CallFunc
	// 接收的调用
	ChanCall chan *CallInfo
}

// NewServer 创建RPC服务器
func NewServer(l int) *Server {
	return &Server{
		functions: make(map[interface{}]CallFunc),
		ChanCall:  make(chan *CallInfo, l),
	}
}

// Register 注册处理函数
// id: reflect.TypeOf(Req)
// f: func(Req) Resp
func (s *Server) Register(id interface{}, f CallFunc) {
	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.functions[id] = f
}

func (s *Server) ret(ci *CallInfo, ri *RetInfo) (err error) {
	if ci.chanRet == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

// Exec 执行调用
func (s *Server) Exec(ci *CallInfo) {
	ret := ci.f(ci.arg)
	err := s.ret(ci, &RetInfo{ret: ret})
	if err != nil {
		log.Error("%v", err)
	}
}

// Go 本地异步调用, 无回调
func (s *Server) Go(id interface{}, arg interface{}) {
	f := s.functions[id]
	if f == nil {
		return
	}

	defer func() {
		recover()
	}()

	s.ChanCall <- &CallInfo{
		f:   f,
		arg: arg,
	}
}

// Call 本地同步调用, 带返回值
func (s *Server) Call(id interface{}, arg interface{}) (interface{}, error) {
	return s.OpenClient(0).Call(id, arg)
}

// Close 关闭rpc服务器并完成所有请求处理
func (s *Server) Close() {
	close(s.ChanCall)

	for ci := range s.ChanCall {
		s.ret(ci, &RetInfo{
			err: errors.New("chanrpc server closed"),
		})
	}
}

// OpenClient 创建客户端
func (s *Server) OpenClient(l int) *Client {
	c := NewClient(l)
	c.Attach(s)
	return c
}
