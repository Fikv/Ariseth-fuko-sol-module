package cache

import (
	"sync"
	"time"
)

type TTLCache[K comparable, V any] struct {
	mu      sync.RWMutex
	ttl     time.Duration
	entries map[K]entry[V]
}

type entry[V any] struct {
	value     V
	expiresAt time.Time
}

func NewTTL[K comparable, V any](ttl time.Duration) *TTLCache[K, V] {
	return &TTLCache[K, V]{
		ttl:     ttl,
		entries: make(map[K]entry[V]),
	}
}

func (c *TTLCache[K, V]) SetTTL(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ttl = ttl
	if ttl <= 0 {
		c.entries = make(map[K]entry[V])
	}
}

func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	var zero V

	c.mu.RLock()
	entry, ok := c.entries[key]
	ttl := c.ttl
	c.mu.RUnlock()

	if !ok || ttl <= 0 {
		return zero, false
	}
	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return zero, false
	}

	return entry.value, true
}

func (c *TTLCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ttl <= 0 {
		return
	}

	c.entries[key] = entry[V]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

func (c *TTLCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[K]entry[V])
}
