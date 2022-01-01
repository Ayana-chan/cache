package lru

import (
	"cache/lru/linkedlist"
)

type Cache struct {
	maxCapacity int //最大容量
	capacity    int //当前容量
	list        *linkedlist.LinkedList
	cache       map[string]*linkedlist.Element
	OnEvicted   func(key string, value linkedlist.Data) //当一个element被删除时执行该方法
}

func New(maxCapacity int, onEvicted func(string, linkedlist.Data)) *Cache {
	return &Cache{
		maxCapacity: maxCapacity,
		list:        linkedlist.New(),
		cache:       make(map[string]*linkedlist.Element),
		OnEvicted:   onEvicted,
	}
}

func (c *Cache) Get(key string) (value linkedlist.Data, ok bool) {
	if element, ok := c.cache[key]; ok {
		list := c.list

		//删除element再将其加到队列头部
		list.Remove(element)
		list.AddToHead(element)

		return element.Value, true
	}
	return linkedlist.Data{}, false
}

func (c *Cache) Add(key string, value linkedlist.Data) {
	if _, ok := c.Get(key); ok {
		c.cache[key].Value = value
	} else {
		element := linkedlist.Element{
			Key:   key,
			Value: value,
		}
		c.list.AddToHead(&element)
		c.cache[key] = &element
		c.capacity += element.Value.Length()

		//容量超过限度
		for c.capacity > c.maxCapacity {
			toRemove := c.list.Tail.Pre
			c.list.Remove(toRemove)
			delete(c.cache, toRemove.Key)
			c.capacity -= toRemove.Value.Length()

			//执行回调方法
			if c.OnEvicted != nil {
				c.OnEvicted(toRemove.Key, toRemove.Value)
			}

		}
	}
}

// GetCapacity 返回当前lru列表的数据量
func (c *Cache) GetCapacity() int {
	return c.capacity
}
