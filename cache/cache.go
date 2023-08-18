package cache

import (
	"sync"
	"time"
)

// example usage:
// var cache = cache.New(24 * time.Hour)

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

func New(dura time.Duration) (cache *Cache) {
	return &Cache{
		m:    make(map[string]item),
		mx:   sync.Mutex{},
		dura: dura,
	}
}

func (ca *Cache) Put(id string, value any) {
	ca.mx.Lock()
	defer ca.mx.Unlock()

	ca.m[id] = item{
		until: time.Now().Add(ca.dura),
		value: value,
	}
}

func (ca *Cache) Get(id string) (value any) {
	ca.mx.Lock()
	defer ca.mx.Unlock()

	item, ok := ca.m[id]
	if !ok {
		return nil
	}

	if item.expired() {
		delete(ca.m, id)
		return nil
	}

	return item
}

func (ca *Cache) Scrub() {
	ca.mx.Lock()
	defer ca.mx.Unlock()

	for id, item := range ca.m {
		if item.expired() {
			delete(ca.m, id)
		}
	}
}
