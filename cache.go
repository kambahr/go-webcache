// Copyright (C) 2022-2026 Kamiar Bahri.
// Use of this source code is governed by
// Boost Software License - Version 1.0
package webcache

import (
	"container/list"
	"sync"
	"time"
)

// CacheItem holds the data and metadata for a cached entry.
type CacheItem struct {
	Path      string
	Content   []byte
	ExpiresAt time.Time
	UserData  map[string]any
}

// entry links the list node back to the map key for efficient eviction.
type entry struct {
	key   string
	value *CacheItem
}

// Cache manages the LRU items with full thread safety.
type Cache struct {
	mu             sync.RWMutex
	maxEntries     int
	defaultTimeout time.Duration
	ll             *list.List
	cache          map[string]*list.Element
}

// NewWebCache creates a new LRU cache instance.
// d: Default expiration duration.
// maxEntries: Max items to keep (0 for unlimited).
func NewWebCache(d time.Duration, maxEntries int) *Cache {
	c := &Cache{
		maxEntries:     maxEntries,
		defaultTimeout: d,
		ll:             list.New(),
		cache:          make(map[string]*list.Element),
	}
	go c.janitor()
	return c
}

// GetItem retrieves content and marks it as "Recently Used".
func (c *Cache) GetItem(path string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ee, ok := c.cache[path]; ok {
		if time.Now().After(ee.Value.(*entry).value.ExpiresAt) {
			c.removeElement(ee)
			return nil, false
		}
		c.ll.MoveToFront(ee)
		return ee.Value.(*entry).value.Content, true
	}
	return nil, false
}

// AddItem adds or updates an item and enforces the LRU limit.
func (c *Cache) AddItem(path string, content []byte, d time.Duration) {
	if d <= 0 {
		d = c.defaultTimeout
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if ee, ok := c.cache[path]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value.Content = content
		ee.Value.(*entry).value.ExpiresAt = time.Now().Add(d)
		return
	}

	item := &CacheItem{
		Path:      path,
		Content:   content,
		ExpiresAt: time.Now().Add(d),
	}
	ele := c.ll.PushFront(&entry{path, item})
	c.cache[path] = ele

	if c.maxEntries > 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

// AddItemDefault adds an item using the default cache duration.
func (c *Cache) AddItemDefault(path string, content []byte) {
	c.AddItem(path, content, c.defaultTimeout)
}

// RemoveItem deletes a specific item by its path.
func (c *Cache) RemoveItem(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ee, ok := c.cache[path]; ok {
		c.removeElement(ee)
	}
}

// ClearAll wipes the entire cache.
func (c *Cache) ClearAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll = list.New()
	c.cache = make(map[string]*list.Element)
}

func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

func (c *Cache) janitor() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for _, ele := range c.cache {
			if now.After(ele.Value.(*entry).value.ExpiresAt) {
				c.removeElement(ele)
			}
		}
		c.mu.Unlock()
	}
}
