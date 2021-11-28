package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    &sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	cache := new(cacheItem)
	cache.key = key
	cache.value = value
	listItem, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(listItem)
		lru.queue.Front().Value = cache
	} else {
		if lru.queue.Len() >= lru.capacity {
			delete(lru.items, lru.queue.Back().Value.(*cacheItem).key)
			lru.queue.Remove(lru.queue.Back())
		}
		lru.items[key] = lru.queue.PushFront(cache)
	}
	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	listItem, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(listItem)
		return listItem.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
