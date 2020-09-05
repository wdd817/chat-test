package chat

import (
	"chat-test/conf"
	"time"
)

type msgLru struct {
	*msg
	next *msgLru
}

// msgLruCache 头尾节点, 单链表
type msgLruCache struct {
	head *msgLru
	foot *msgLru
}

func newMsgLruCache() *msgLruCache {
	cache := &msgLruCache{}
	return cache
}

func (cache *msgLruCache) Add(m *msg) {
	newLru := &msgLru{ msg: m }
	if cache.foot == nil {
		cache.head = newLru
		cache.foot = newLru
	} else {
		cache.foot.next = newLru
		cache.foot = newLru
	}

	for cache.head.timestamp < time.Now().Unix() - int64(conf.PopularSecond) {
		cache.head = cache.head.next
		if cache.head == nil {
			cache.foot = nil
			break
		}
	}
}

func (cache *msgLruCache) Popular() string {
	var maxCnt int
	var popular string
	allCnt := make(map[string]int)
	for cur := cache.head; cur != nil; cur = cur.next {
		cnt := allCnt[cur.content] + 1
		if cnt > maxCnt {
			maxCnt = cnt
			popular = cur.content
		}
	}

	return popular
}
