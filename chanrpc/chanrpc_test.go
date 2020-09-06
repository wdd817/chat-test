package chanrpc

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestChanRPC chanrpc 测试用例
func TestChanRPC(t *testing.T) {
	s := NewServer(10)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		s.Register("toString", func(req interface{}) interface{} {
			return fmt.Sprint(req)
		})

		wg.Done()

		for {
			s.Exec(<-s.ChanCall)
		}
	}()

	wg.Wait()
	wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic with recover: %v", r)
			}
		}()

		c := s.OpenClient(10)
		ret, err := c.Call("toString", 1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		retStr, ok := ret.(string)
		if !ok {
			t.Errorf("unexpected toString ret value type")
		}

		if retStr != fmt.Sprint(1) {
			t.Errorf("unexpected call ret: %s", retStr)
		}

		c.Async("toString", 1, func(ret interface{}, err error) {
			if err != nil {
				t.Errorf("unexpected async call error: %v", err)
			}

			if retStr != fmt.Sprint(1) {
				t.Errorf("unexpected async ret: %s", retStr)
			}
		})

		c.Cb(<-c.ChanAsyncRet)

		s.Go("toString", 1)

		wg.Done()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
		t.Errorf("unexpected async timeout")
	}()

	wg.Wait()
}
