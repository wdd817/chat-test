package chat

import (
	"testing"
)

func makeCircleData() *msgCircleCache {
	cache := newMsgCircleCache(3)
	cache.Add(&msg{
		sender: "tester",
		content: "hello",
		timestamp: 1,
	})
	cache.Add(&msg{
		sender: "tester",
		content: "world",
		timestamp: 2,
	})
	cache.Add(&msg{
		sender: "tester",
		content: "foo",
		timestamp: 3,
	})
	cache.Add(&msg{
		sender: "tester",
		content: "bar",
		timestamp: 4,
	})
	cache.Add(&msg{
		sender: "tester",
		content: "baz",
		timestamp: 5,
	})
	return cache
}

func TestMsgCircleCache(t *testing.T) {
	cache := makeCircleData()

	var i int
	cache.Range(func(m *msg) {
		if i == 0 && m.content != "foo" {
			t.Errorf("expected %s got %s", "foo", m.content)
		}
		if i == 1 && m.content != "bar" {
			t.Errorf("expected %s got %s", "bar", m.content)
		}
		if i == 2 && m.content != "baz" {
			t.Errorf("expected %s got %s", "baz", m.content)
		}
		if i > 2 {
			t.Error("too much cache")
		}
		i++
	})
}

