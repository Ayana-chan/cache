package cachestruct

import (
	"cache/lru"
	"sync"
)

// Cache 封装了lru包中的cache，为其添加并发特性，并以container为存储的数据
type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	CacheBytes int64 //lru的最大容量
}

func (c *Cache) add(key string, value Container) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.CacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) get(key string) (value Container, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(Container), ok
	}
	return
}
