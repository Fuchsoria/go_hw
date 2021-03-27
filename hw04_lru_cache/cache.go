package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
	Len() int
	QueueLen() int
}

type lruCache struct {
	sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	var alreadyExist bool

	if prevValue, ok := c.items[key]; ok {
		c.queue.Remove(prevValue)

		alreadyExist = true
	} else {
		alreadyExist = false

		c.checkCapacity(1)
	}

	c.items[key] = c.queue.PushFront(value)
	c.Unlock()

	return alreadyExist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()

	if item, ok := c.items[key]; ok {
		c.Unlock()

		return item.Value, true
	}

	c.Unlock()

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue.ClearFrontBack()
	c.items = make(map[Key]*ListItem)
}

func (c *lruCache) checkCapacity(extra int) {
	if c.capacity < (c.queue.Len() + extra) {
		l := c.queue

		l.Remove(l.Back())
	}
}

func (c *lruCache) Len() int {
	return len(c.items)
}

func (c *lruCache) QueueLen() int {
	return c.queue.Len()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
