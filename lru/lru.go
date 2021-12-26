package lru

import (
	"cache/lru/linkedlist"
)

type Cache struct {
	maxCapacity int64
	capacity    int64
	list        *linkedlist.LinkedList
	cache       map[string]*linkedlist.Element
	OnEvicted   func(key string, value linkedlist.Value)
}

func New(maxCapacity int64, onEvicted func(string, linkedlist.Value)) *Cache {
	return &Cache{
		maxCapacity: maxCapacity,
		list:        linkedlist.New(),
		cache:       make(map[string]*linkedlist.Element),
		OnEvicted:   onEvicted,
	}
}

func (c *Cache) Get(key string) (value linkedlist.Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		list := c.list
		list.Remove(element)
		list.AddToHead(element)
		return element.Value, true
	}
	return nil, false
}

func (c *Cache) Add(key string, value linkedlist.Value) {
	if _, ok := c.Get(key); ok {
		c.cache[key].Value = value
	} else {
		element := linkedlist.Element{
			Key:   key,
			Value: value,
		}
		if c.capacity == c.maxCapacity {
			toRemove := c.list.Tail.Pre
			c.list.Remove(toRemove)
			delete(c.cache, toRemove.Key)
			c.capacity--
			if c.OnEvicted != nil {
				c.OnEvicted(toRemove.Key, toRemove.Value)
			}
		}
		c.list.AddToHead(&element)
		c.cache[key] = &element
		c.capacity++
	}
}

func (c *Cache) GetCapacity() int64 {
	return c.capacity
}
