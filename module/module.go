package module

import (
	"chat-test/chanrpc"
)

type Module interface {
	OnInit() error
	OnDestroy()
	Run(closeSig chan bool)
	Name() string
	ChanRPC() *chanrpc.Server
}
