package main

import (
	"log"
	"sync"
	"time"
)

type timedElement struct {
	v string
	t time.Time
}

// SafeCache is safe to use concurrently.
type SafeCache struct {
	elements map[string]*timedElement
	limit    time.Duration
	mutex    sync.Mutex
}

// NewSafeCache creates a new SafeCache and starts a goroutine to remove elements after seconds "limit".
func NewSafeCache(limit int) *SafeCache {
	cache := &SafeCache{
		elements: make(map[string]*timedElement),
		limit:    time.Duration(limit) * time.Second,
	}

	go timeDel(time.Tick(1*time.Second), cache)

	return cache
}

// Add sets a key "k" to a value "v".
func (c *SafeCache) Add(k, v string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.elements[k] = &timedElement{v: v, t: time.Now()}
}

// Get returns a value stored under a key "k".
func (c *SafeCache) Get(k string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[k]; ok {
		elem.t = time.Now()
		return elem.v
	}

	return ""
}

// Delete removes an entry from the cache
func (c *SafeCache) Delete(k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.elements, k)
}

func timeDel(ticker <-chan time.Time, cache *SafeCache) {
	for range ticker {
		for k, elem := range cache.elements {
			if time.Since(elem.t) >= cache.limit {
				log.Printf("deleting %s", k)
				cache.Delete(k)
			}
		}
	}
}
