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
	mutex    sync.Mutex
}

type queueItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.items[key]
	if !ok {
		if c.len() >= c.capacity {
			delete(c.items, c.queue.Back().Value.(queueItem).Key)
			c.queue.Remove(c.queue.Back())
		}
		c.queue.PushFront(queueItem{key, value})
		c.items[key] = c.queue.Front()
	} else {
		c.items[key].Value = queueItem{key, value}
		c.queue.MoveToFront(c.items[key])
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	result := item.Value.(queueItem).Value
	return result, true
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) len() int {
	return c.queue.Len()
}
