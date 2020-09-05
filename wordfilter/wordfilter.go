package wordfilter

import (
	"chat-test/chanrpc"
	"chat-test/conf"
	"chat-test/log"
	"chat-test/module"
	"chat-test/proto"
	"fmt"
	"io/ioutil"
	"strings"
)

const profanityWordsFile = "conf/profanity-words.txt"

type WordFilter struct {
	skeleton *module.Skeleton

	dfaUtil *DFAUtil
}

func (w WordFilter) Run(closeSig chan bool) {
	w.skeleton.Run(closeSig)
	log.Info("world filter stopping")
}

func NewModule() module.Module {
	wordFilter := &WordFilter{
		skeleton: &module.Skeleton{
			AsynCallLen:   conf.AsynCallLen,
			ChanRPCServer: chanrpc.NewServer(conf.ChanRPCLen),
		},
	}
	wordFilter.skeleton.Init()
	return wordFilter
}

func (w WordFilter) OnInit() error {
	bytes, err := ioutil.ReadFile(profanityWordsFile)
	if err != nil {
		return fmt.Errorf("cannot read profanity words file: %v", err)
	}

	words := strings.Split(string(bytes), "\n")
	w.dfaUtil = NewDFAUtil(words)
	w.skeleton.RegisterChanRPC((*proto.WordFilterReq)(nil), w.OnWordFilterReq)
	return nil
}

func (w WordFilter) OnDestroy() {
	log.Info("world filter stopped")
}

func (w WordFilter) Name() string {
	return conf.WordFilterModule
}

func (w WordFilter) ChanRPC() *chanrpc.Server {
	return w.skeleton.ChanRPCServer
}

func (w *WordFilter) OnWordFilterReq(req0 interface{}) interface{} {
	req := req0.(*proto.WordFilterReq)
	filtered, isFiltered := w.dfaUtil.HandleWord(req.Content, '*')
	return &proto.WordFilterResp{
		IsFiltered: isFiltered,
		Filtered:   filtered,
	}
}
