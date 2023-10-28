package hw04lrucache

import (
	"log"
	"sync"
)

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

type Value struct {
	key   Key
	value interface{}
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

	prepValue := Value{
		key:   key,
		value: value,
	}

	switch ok {
	case true:
		val.Value = prepValue
		c.queue.MoveToFront(val)
	case false:
		newItem := c.queue.PushFront(prepValue)
		c.items[key] = newItem
		if len(c.items) > c.capacity {
			back := c.queue.Back()
			v, successfully := back.Value.(Value)
			if !successfully {
				log.Println("setting error: unsupported value type")
			}
			delete(c.items, v.key)
			c.queue.Remove(back)
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
		switch v := item.Value.(type) {
		case Value:
			val = v.value
		default:
			log.Println("getting error: unsupported value type")
			val = nil
		}

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
