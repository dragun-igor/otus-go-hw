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
	lru.mutex.Lock() // Mutex для потокобезопасности
	defer lru.mutex.Unlock() // Так делать нехорошо, но конкретной задачи нету
	// Создаём новый кэш айтем, добавляем в него ключи и значение
	// Если ключ существовал, то переносим в начало очереди и обновляем значение
	// Если ключа не было, то добавляем значения в конец очереди
	// Если очередь была полная - удаляем последний элемент
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
	// Всё просто, если запрашиваем значение по ключу, если ключ существовал, то переносим в начало очереди
	listItem, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(listItem)
		return listItem.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (lru *lruCache) Clear() {
	// Создаём пустой список и привязываем к очереди
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
