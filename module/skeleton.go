package module

import "chat-test/chanrpc"

type Skeleton struct {
	ChanRPCServer *chanrpc.Server
	server        *chanrpc.Server
	client        *chanrpc.Client
}

func (s *Skeleton) Init() {
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
		case ri := <- s.client.ChanRet:
			s.client.Cb(ri)
		case ci := <- s.server.ChanCall:
			s.server.Exec(ci)
		}
	}
}

func (s *Skeleton) Async(server *chanrpc.Server, id interface{}, arg interface{}) {
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f chanrpc.CallFunc) {
	if s.ChanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.server.Register(id, f)
}
