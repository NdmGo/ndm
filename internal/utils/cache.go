package utils

import (
	"sync"
	"time"
)

// CacheItem cache item
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
	Accessed   time.Time
}

// IsExpired check if expired
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.Expiration)
}

// MemoryCache memory cache
type MemoryCache struct {
	items      map[string]*CacheItem
	mu         sync.RWMutex
	maxSize    int
	defaultTTL time.Duration
	stopCh     chan struct{}
	started    bool
}

// NewMemoryCache create new memory cache
func NewMemoryCache(maxSize int, defaultTTL time.Duration) *MemoryCache {
	if maxSize <= 0 {
		maxSize = 1000 // default maximum cache size
	}
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute // default TTL
	}

	cache := &MemoryCache{
		items:      make(map[string]*CacheItem),
		maxSize:    maxSize,
		defaultTTL: defaultTTL,
		stopCh:     make(chan struct{}),
	}

	// start cleanup goroutine
	go cache.cleanup()
	cache.started = true

	return cache
}

// Set set cache item
func (c *MemoryCache) Set(key string, value interface{}, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if cache is full, delete the oldest item
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	expiration := c.defaultTTL
	if len(ttl) > 0 && ttl[0] > 0 {
		expiration = ttl[0]
	}

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
		Accessed:   time.Now(),
	}
}

// Get get cache item
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists || item.IsExpired() {
		if exists {
			c.Delete(key)
		}
		return nil, false
	}

	// update access time
	c.mu.Lock()
	item.Accessed = time.Now()
	c.mu.Unlock()

	return item.Value, true
}

// Delete delete cache item
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear clear cache
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Size get cache size
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys get all keys
func (c *MemoryCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// GetStats get cache statistics
func (c *MemoryCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiredCount := 0
	for _, item := range c.items {
		if item.IsExpired() {
			expiredCount++
		}
	}

	return CacheStats{
		TotalItems:   len(c.items),
		ExpiredItems: expiredCount,
		MaxSize:      c.maxSize,
		HitRate:      0, // requires additional statistics
	}
}

// CacheStats cache statistics
type CacheStats struct {
	TotalItems   int     `json:"total_items"`
	ExpiredItems int     `json:"expired_items"`
	MaxSize      int     `json:"max_size"`
	HitRate      float64 `json:"hit_rate"`
}

// evictOldest delete the oldest cache item
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestKey == "" || item.Accessed.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.Accessed
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

// cleanup clean up expired items
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.items {
				if item.IsExpired() {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

// Stop stop cache cleanup
func (c *MemoryCache) Stop() {
	if c.started {
		close(c.stopCh)
		c.started = false
	}
}

// 全局缓存实例
var (
	DefaultCache *MemoryCache
	cacheOnce    sync.Once
)

// GetDefaultCache 获取默认缓存
func GetDefaultCache() *MemoryCache {
	cacheOnce.Do(func() {
		DefaultCache = NewMemoryCache(0, 0)
	})
	return DefaultCache
}

// 便捷函数
func CacheSet(key string, value interface{}, ttl ...time.Duration) {
	GetDefaultCache().Set(key, value, ttl...)
}

func CacheGet(key string) (interface{}, bool) {
	return GetDefaultCache().Get(key)
}

func CacheDelete(key string) {
	GetDefaultCache().Delete(key)
}

func CacheClear() {
	GetDefaultCache().Clear()
}
