package webcache

import (
	"strings"
	"time"
)

// CacheItem holds the attributes of a request file to be
type CacheItem struct {
	Path            string
	Content         []byte
	Expires         time.Duration
	UserData        map[string]interface{}
	DateTimeCreated time.Time
}

var mCacheArry []CacheItem
var cacheDirty bool

// Cache type holds the global required attibutes.
type Cache struct {
	CacheDuration time.Duration
}

// NewWebCache creates a new instance of webCache.
// Full root physical path of the website, default cache duration
func NewWebCache(d time.Duration) *Cache {
	var c Cache
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
// expired items. A goto jump is used to avoid recursion.
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

	time.Sleep(800 * time.Millisecond)

	goto lblAgain
}

//-----------------------------------------------
//               public functions               '
//-----------------------------------------------

// ClearAll removes all cache.
func (c *Cache) ClearAll() {
	mCacheArry = make([]CacheItem, 0)
}

// Clears an items from the global bache.
func (c *Cache) Clear(path string) {
	for i := 0; i < len(mCacheArry); i++ {
		if strings.HasPrefix(mCacheArry[i].Path, path) {
			removeItem(i)
			return
		}
	}
}

// GetCacheList makes the mCacheArry visible.
func (c *Cache) GetCacheList(uriPath string) []CacheItem {
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

// GetItemDetailed returns a selected item from the global array.
func (c *Cache) GetItemDetailed(uriPath string) CacheItem {
	var b CacheItem
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == uriPath {
			return mCacheArry[i]
		}
	}

	return b
}

// GetItem returns a selected item from the global array.
func (c *Cache) GetItem(uriPath string) []byte {
	var b []byte
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == uriPath {
			return mCacheArry[i].Content
		}
	}

	return b
}

// AddItemDetailed add a cache item to the global list with the
// added user data.
func (c *Cache) AddItemDetailed(uriPath string, content []byte, d time.Duration, userData map[string]interface{}) {
	var cx CacheItem
	cx.Path = uriPath
	cx.Content = content
	cx.DateTimeCreated = time.Now()
	cx.Expires = d
	cx.UserData = userData
	mCacheArry = append(mCacheArry, cx)

	if !cacheDirty {
		go c.manageCache()
		cacheDirty = true
	}
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

// AddItem adds an item to the global array.
func (c *Cache) AddItemDefault(uriPath string, content []byte) {
	c.AddItem(uriPath, content, c.CacheDuration)
}

// RemoveItem removes an item from the global array.
func (c *Cache) RemoveItem(p string) {
	for i := 0; i < len(mCacheArry); i++ {
		if mCacheArry[i].Path == p {
			remove(mCacheArry, i)
			break
		}
	}
}
