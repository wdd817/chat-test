package module

import (
	"chat-test/chanrpc"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

type testReq struct {
	val int
}

type testResp struct {
	val string
}

// TestSkeleton skeleton 测试用例
func TestSkeleton(t *testing.T) {
	s1 := Skeleton{
		AsynCallLen:   10,
		ChanRPCServer: chanrpc.NewServer(10),
	}
	s2 := Skeleton{
		AsynCallLen:   10,
		ChanRPCServer: chanrpc.NewServer(10),
	}
	s1.Init()
	s2.Init()

	closeSig := make(chan bool)

	go func() {
		s1.Run(closeSig)
	}()
	go func() {
		s2.Run(closeSig)
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		toString := func(req0 interface{}) interface{} {
			req, _ := req0.(*testReq)
			return &testResp{val: fmt.Sprint(req.val)}
		}
		s2.RegisterChanRPC((*testReq)(nil), toString)


		s1.Async(s2.ChanRPCServer, &testReq{val: 1}, func(ret interface{}, err error) {
			defer wg.Done()

			if err != nil {
				t.Errorf("unexpected s1 async err: %v", err)
				return
			}

			resp, ok := ret.(*testResp)
			if !ok {
				t.Errorf("unexpected s1 async ret type: %v", reflect.TypeOf(resp))
				return
			}

			if resp.val != "1" {
				t.Errorf("unexpected s1 async ret: %s", resp.val)
			}
		})
	}()

	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
		t.Errorf("unexpected async timeout")
	}()

	wg.Wait()
	closeSig <- true
}
