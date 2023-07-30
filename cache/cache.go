package cache

import (
	"sync"
	"time"
)

type m map[string]item
type Cache struct {
	m
	mx   sync.Mutex
	dura time.Duration
}

type item struct {
	until time.Time
	value any
}

func (item item) expired() bool {
	return item.until.After(time.Now())
}

func New(duration time.Duration) (cache *Cache) {
	cache.m = make(m)
	cache.dura = duration
	return cache
}

func (cache *Cache) Put(id string, value any) {
	cache.mx.Lock()
	defer cache.mx.Unlock()

	cache.m[id] = item{
		until: time.Now().Add(cache.dura),
		value: value,
	}
}

func (cache *Cache) Get(id string) (value any) {
	cache.mx.Lock()
	defer cache.mx.Unlock()

	item, ok := cache.m[id]
	if !ok {
		return nil
	}

	if item.expired() {
		delete(cache.m, id)
		return nil
	}

	return item
}
