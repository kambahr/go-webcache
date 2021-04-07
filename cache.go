package webcache

import (
	"time"
)

// CacheItem holds the attributes of a request file to be
type CacheItem struct {
	Path            string
	Content         []byte
	Expires         time.Duration
	DateTimeCreated time.Time
}

var mCacheArry []CacheItem
var cacheDirty bool

// Cache type holds the global required attibutes.
type Cache struct {
	RootPhysicalPath string
	CacheDuration    time.Duration
}

// NewWebCache creates a new instance of webCache.
// Full root physical path of the website, default cache duration
func NewWebCache(rootPhysicalPath string, d time.Duration) *Cache {
	var c Cache
	c.RootPhysicalPath = rootPhysicalPath
	c.CacheDuration = d
	return &c
}

// remove drops an item from the cache list.
func remove(s []CacheItem, i int) []CacheItem {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// removeItem removes an item from the global cache array by its index position.
func removeItem(idx int) {
	for i := 0; i < len(mCacheArry); i++ {
		if i == idx {
			mCacheArry = remove(mCacheArry, i)
			return
		}
	}
}

// manageCache goes through the global cache array and removes the
// expired items -- every 5 seconds. A goto jump is used to avoid recursion,
// which could lead to memory isues for long running processes.
func (c *Cache) manageCache() {
lblAgain:
	for i := 0; i < len(mCacheArry); i++ {

		elapsed := time.Since(mCacheArry[i].DateTimeCreated)

		// Apply the value set specifically for this file - first.
		if elapsed >= mCacheArry[i].Expires {
			removeItem(i)
			break
		}

		// Apply the global cache duration,
		if elapsed >= c.CacheDuration {
			removeItem(i)
			break
		}
	}

	time.Sleep(750 * time.Millisecond)

	goto lblAgain
}

//-----------------------------------------------
//               public functions               '
//-----------------------------------------------

// GetCacheList makes the mCacheArry visible.
func (*Cache) GetCacheList(uriPath string) []CacheItem {
	return mCacheArry
}

// Exists tells if an item exists.
func (c *Cache) Exists(uriPath string) bool {
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == uriPath {
			return true
		}
	}
	return false
}

// GetItem returns a selected item from the global array.
func (*Cache) GetItem(uriPath string) []byte {
	var b []byte
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == uriPath {
			return mCacheArry[i].Content
		}
	}

	return b
}

// AddItem adds an item to the global array.
func (c *Cache) AddItem(uriPath string, content []byte, d time.Duration) {
	var cx CacheItem
	cx.Path = uriPath
	cx.Content = content
	cx.DateTimeCreated = time.Now()
	cx.Expires = d
	mCacheArry = append(mCacheArry, cx)

	if !cacheDirty {
		go c.manageCache()
		cacheDirty = true
	}
}

// RemoveItem removes an item from the global array.
func (*Cache) RemoveItem(p string) {
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == p {
			remove(mCacheArry, i)
			break
		}
	}
}
