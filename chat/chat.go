package chat

import (
	"chat-test/chanrpc"
	"chat-test/conf"
	"chat-test/log"
	"chat-test/module"
	"chat-test/proto"
)

type Chat struct {
	skeleton *module.Skeleton

	history *msgCircleCache
	recent  *msgLruCache

	userIDs map[string]int
	users   map[int]*User
}

func NewModule() module.Module {
	chat := &Chat{
		skeleton: &module.Skeleton{
			AsynCallLen:   conf.AsynCallLen,
			ChanRPCServer: chanrpc.NewServer(conf.ChanRPCLen),
		},
		history: newMsgCircleCache(conf.MaxHistory),
		recent:  newMsgLruCache(),
		userIDs: make(map[string]int),
		users:   make(map[int]*User),
	}
	chat.skeleton.Init()
	return chat
}

func (chat *Chat) OnInit() error {

	chat.skeleton.RegisterChanRPC((*proto.LoginReq)(nil), chat.OnLoginReq)
	chat.skeleton.RegisterChanRPC((*proto.ChatMsgReq)(nil), chat.OnChatMsgReq)
	chat.skeleton.RegisterChanRPC((*proto.PopularReq)(nil), chat.OnPopularReq)
	chat.skeleton.RegisterChanRPC((*proto.StatsReq)(nil), chat.OnStatsReq)
	chat.skeleton.RegisterChanRPC((*proto.LogoutReq)(nil), chat.OnLogoutReq)

	return nil
}

func (chat *Chat) OnDestroy() {
	log.Info("chat stopped")
}

func (chat *Chat) Run(closeSig chan bool) {
	chat.skeleton.Run(closeSig)
	log.Info("chat stopping")
}

func (chat *Chat) Name() string {
	return conf.ChatModule
}

func (chat *Chat) ChanRPC() *chanrpc.Server {
	return chat.skeleton.ChanRPCServer
}
