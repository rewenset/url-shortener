package main

import (
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
	mutex    sync.Mutex
}

// NewSafeCache creates a new SafeCache.
func NewSafeCache() *SafeCache {
	return &SafeCache{
		elements: make(map[string]*timedElement),
	}
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

	if element, ok := c.elements[k]; ok {
		return element.v
	}

	return ""
}

// Delete removes an entry from the cache
func (c *SafeCache) Delete(k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.elements, k)
}
