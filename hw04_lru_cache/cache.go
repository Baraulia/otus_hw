package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.items[key]
	if ok {
		val.Value = value
		c.queue.Front()
	} else {
		newItem := c.queue.PushFront(value)
		c.items[key] = newItem
		if len(c.items) > c.capacity {
			back := c.queue.Back()
			c.queue.Remove(back)
			c.deleteFromMap(back)
		}
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var val interface{}
	item, ok := c.items[key]
	if ok {
		val = item.Value
		c.queue.PushFront(item)
	}

	return val, ok
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func (c *lruCache) deleteFromMap(value *ListItem) {
	for k, v := range c.items {
		if v == value {
			delete(c.items, k)
		}
	}
}
