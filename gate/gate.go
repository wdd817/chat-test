package gate

import (
	"chat-test/chanrpc"
	"chat-test/conf"
	"chat-test/log"
	"chat-test/module"
	"chat-test/network"
	"chat-test/network/text"
	"chat-test/proto"
)

type Gateway struct {
	skeleton *module.Skeleton

	wss       *network.WSServer
	Processor network.Processor

	agents map[int]*agent
}

func NewModule() module.Module {
	return newGateway()
}

func newGateway() *Gateway {
	gateway := &Gateway{
		skeleton: &module.Skeleton{
			AsynCallLen:   conf.AsynCallLen,
			ChanRPCServer: chanrpc.NewServer(conf.ChanRPCLen),
		},
		agents: make(map[int]*agent),
	}
	gateway.skeleton.Init()
	return gateway
}

func (gate *Gateway) OnInit() error {
	processor := text.NewProcessor()
	processor.SetDefaultHandler(gate.handleChat)
	processor.SetHandler("/popular", gate.handlePopular)
	processor.SetHandler("/stats", gate.handleStats)
	gate.Processor = processor

	gate.skeleton.RegisterChanRPC((*proto.BroadcastReq)(nil), gate.OnBroadcastReq)

	gate.wss = &network.WSServer{
		Addr:            conf.GatewayAddr,
		MaxConnNum:      conf.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		HTTPTimeout:     conf.HTTPTimeout,
		NewAgent: func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			a.WriteMsg("Please enter your name: ")
			return a
		},
	}

	gate.wss.Start()

	return nil
}

func (gate *Gateway) OnDestroy() {
	log.Info("gateway stopped")
}

func (gate *Gateway) Name() string {
	return conf.GateModule
}

func (gate *Gateway) ChanRPC() *chanrpc.Server {
	return gate.skeleton.ChanRPCServer
}

func (gate *Gateway) Run(closeSig chan bool) {
	gate.skeleton.Run(closeSig)
	log.Info("gateway stopping")
}
