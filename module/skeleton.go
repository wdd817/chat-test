package module

import (
	"chat-test/chanrpc"
	"chat-test/log"
	"reflect"
)

type Skeleton struct {
	AsynCallLen        int
	ChanRPCServer *chanrpc.Server
	server        *chanrpc.Server
	client        *chanrpc.Client
}

func (s *Skeleton) Init() {
	if s.AsynCallLen <= 0 {
		s.AsynCallLen = 0
	}

	s.client = chanrpc.NewClient(s.AsynCallLen)
	s.server = s.ChanRPCServer

	if s.server == nil {
		s.server = chanrpc.NewServer(0)
	}
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			s.server.Close()
			for !s.client.Idle() {
				s.client.Close()
			}
			return
		case ri := <- s.client.ChanAsyncRet:
			s.client.Cb(ri)
		case ci := <- s.server.ChanCall:
			s.server.Exec(ci)
		}
	}
}

func (s *Skeleton) Async(server *chanrpc.Server, arg interface{}, cb chanrpc.Callback) {
	if s.AsynCallLen == 0 {
		log.Error("invalid AsyncCallLen")
		return
	}

	idType := reflect.TypeOf(arg)
	if idType == nil || idType.Kind() != reflect.Ptr {
		log.Fatal("not use pointer to call: %v %v", arg, idType)
		return
	}

	s.client.Attach(server)
	s.client.Async(idType.Elem().Name(), arg, cb)
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f chanrpc.CallFunc) {
	if s.ChanRPCServer == nil {
		log.Fatal("invalid ChanRPCServer")
		return
	}

	idType := reflect.TypeOf(id)
	if idType == nil || idType.Kind() != reflect.Ptr {
		log.Fatal("not use nil type pointer to register: %v %v", id, idType)
		return
	}

	s.server.Register(idType.Elem().Name(), f)
}
