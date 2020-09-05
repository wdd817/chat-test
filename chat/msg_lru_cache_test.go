package chat

import (
	"chat-test/conf"
	"testing"
	"time"
)

func TestFoo(t *testing.T) {
	h := newMsgLruCache()
	now := time.Now().Unix()
	h.Add(&msg{
		content:   "hello",
		timestamp: now - int64(conf.PopularSecond) - 3,
	})
	h.Add(&msg{
		content:   "world",
		timestamp: now - int64(conf.PopularSecond) - 2,
	})
	h.Add(&msg{
		content:   "foo",
		timestamp: now - int64(conf.PopularSecond) - 1,
	})
	h.Add(&msg{
		content:   "bar",
		timestamp: now - int64(conf.PopularSecond),
	})
	h.Add(&msg{
		content:   "baz",
		timestamp: now - int64(conf.PopularSecond) + 1,
	})

	if h.head.content != "bar" {
		t.Errorf("expected %s got %s", "bar", h.head.content)
	}

	if h.foot.content != "baz" {
		t.Errorf("expected %s got %s", "baz", h.foot.content)
	}

	time.Sleep(time.Duration(conf.PopularSecond) * time.Second)
	h.Add(&msg{
		content: "expire",
		timestamp: 0,
	})
	if h.head != nil || h.foot != nil {
		t.Errorf("list should be empty but got %v %v", h.head, h.foot)
	}
}
