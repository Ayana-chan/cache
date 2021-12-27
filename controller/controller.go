package controller

import (
	"cache/cachestruct"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 接口型函数，可在使用中代替Getter
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cachestruct.Cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cachestruct.Cache{CacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (cachestruct.Data, error) {
	if key == "" {
		return cachestruct.Data{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("Cache Hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value cachestruct.Data, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (cachestruct.Data, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return cachestruct.Data{}, err

	}
	value := cachestruct.Data{B: cachestruct.CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value cachestruct.Data) {
	g.mainCache.Add(key, value)
}
