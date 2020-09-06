package app

import (
	"chat-test/chanrpc"
	"chat-test/module"
	"sync"
	"testing"
	"time"
)

type testModule struct {
	*chanrpc.Server
}

func (t testModule) OnInit() error            { return nil }
func (t testModule) OnDestroy()               {}
func (t testModule) Run(closeSig chan bool)   { <-closeSig }
func (t testModule) Name() string             { return "test" }
func (t testModule) ChanRPC() *chanrpc.Server { return nil }

func newTestModule() module.Module {
	return &testModule{
		Server: chanrpc.NewServer(1),
	}
}

// TestApp app测试用例
func TestApp(t *testing.T) {
	// 创建module
	testModule := newTestModule()

	if Instance().GetState() != StateNone {
		t.Errorf("error instance")
		return
	}

	// 运行app, 然后等待1秒
	go Instance().Run(testModule)
	time.Sleep(1 * time.Second)
	if Instance().GetState() != StateRun {
		t.Errorf("error starting")
		return
	}

	rpc := Instance().GetChanRPC(testModule.Name())
	if rpc != testModule.ChanRPC() {
		t.Errorf("find different chan rpc expect: %v got %v", rpc, testModule.ChanRPC())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// 关闭超时
	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	// 停止app
	go func() {
		Instance().Terminate()
		wg.Done()
	}()

	wg.Wait()
	if Instance().GetState() != StateNone {
		t.Errorf("error termating")
		return
	}
}
