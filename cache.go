package geeCache

import (
	"geeCache/lru"
	"sync"
)

//并发控制
type cache struct {
	mu         sync.Mutex
	lru        lru.Cache
	cacheBytes int64
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val, ok := c.lru.Get(key); ok {
		return val.(ByteView), ok
	}
	return
}

func (c *cache) Add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Add(key, val)
}
