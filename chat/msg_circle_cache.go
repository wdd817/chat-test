package chat

import (
	"chat-test/log"
)

type msgCircleCache struct {
	maxLen int
	data   []*msg
	curIdx int
}

func newMsgCircleCache(maxLen int) *msgCircleCache {
	if maxLen <= 0 {
		log.Error("newMsgCircleCache maxLen less than 1: %d", maxLen)
		return nil
	}

	cache := &msgCircleCache{
		maxLen: maxLen,
		data:   make([]*msg, maxLen),
	}
	return cache
}

func (cache *msgCircleCache) Add(m *msg) {
	cache.data[cache.curIdx] = m
	cache.curIdx++
	if cache.curIdx == cache.maxLen {
		cache.curIdx = 0
	}
}

func (cache *msgCircleCache) Range(f func(m *msg)) {
	idx := cache.curIdx
	endIdx := cache.curIdx - 1
	if endIdx < 0 {
		endIdx = cache.maxLen - 1
	}
	for idx != endIdx {
		m := cache.data[idx]
		if m != nil {
			f(m)
		}
		idx = (idx + 1) % cache.maxLen
	}
}
