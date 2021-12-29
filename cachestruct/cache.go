package cachestruct

import (
	"cache/lru"
	"cache/lru/linkedlist"
	"sync"
)

// Cache 封装了lru包中的cache，为其添加并发特性
type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	CacheBytes int64 //lru的最大容量,单位是字节
}

func (c *Cache) Add(key string, value linkedlist.Data) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {

		//延迟初始化lru
		c.lru = lru.New(c.CacheBytes, nil)

	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value linkedlist.Data, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v, ok
	}
	return
}
