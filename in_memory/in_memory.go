package in_memory

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
	"time"
)

// CacheNode represents a single node in the LRU cache with a key-value pair and expiration time.
type CacheNode struct {
	key      string    // Key of the cache entry
	value    string    // Value associated with the key
	expireAt time.Time // Expiration time for the cache entry (zero time if no expiration)
}

// LRUCache implements a Least Recently Used (LRU) cache using a map and a doubly linked list.
type LRUCache struct {
	cache map[string]*list.Element // Map for fast access to cache elements
	list  *list.List               // Doubly linked list to track access order

	cleanupTime time.Duration // Time interval for periodic cleanup of expired entries
	mu          sync.Mutex    // Mutex for concurrent access to cache data structures
}

// NewLRUCache initializes and returns a new LRUCache instance.
func NewLRUCache(cleanupTime time.Duration) *LRUCache {
	c := &LRUCache{
		cache: make(map[string]*list.Element),
		list:  list.New(),

		cleanupTime: cleanupTime,
	}
	go c.startCleanupRoutine() // Start a goroutine for periodic cache cleanup
	return c
}

// startCleanupRoutine starts a background goroutine to clean up expired items from the cache periodically.
func (c *LRUCache) startCleanupRoutine() {
	ticker := time.NewTicker(c.cleanupTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup removes expired items from the cache.
func (c *LRUCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for elem := c.list.Front(); elem != nil; {
		next := elem.Next()
		node := elem.Value.(*CacheNode)
		if !node.expireAt.IsZero() && node.expireAt.Before(now) {
			// Remove expired node from the linked list and delete from map
			c.list.Remove(elem)
			delete(c.cache, node.key)
		}
		elem = next
	}
}

// Get retrieves the value associated with the given key.
// It moves the accessed element to the front of the list to mark it as recently used.
func (c *LRUCache) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.cache[key]; found {
		node := elem.Value.(*CacheNode)
		if node.expireAt.IsZero() || node.expireAt.After(time.Now()) {
			c.list.MoveToFront(elem) // Move accessed item to the front of the list
			return node.value
		}
		// Remove the expired element from both the list and the map
		c.list.Remove(elem)
		delete(c.cache, key)
	}
	return "" // Return empty string if key not found or expired
}

// Put adds a key-value pair to the cache with an optional TTL.
// If the key already exists, it updates the value and moves the element to the front.
// If the cache exceeds maxLength, it evicts the least recently used element.
func (c *LRUCache) Put(key string, value string, length int, ttl int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.cache[key]; found {
		node := elem.Value.(*CacheNode)
		node.value = value
		if ttl > 0 {
			node.expireAt = time.Now().Add(time.Duration(ttl) * time.Second)
		} else {
			node.expireAt = time.Time{} // Reset expiration if ttl <= 0
		}
		c.list.MoveToFront(elem) // Move existing item to the front
		return
	}
	if c.list.Len() >= length {
		c.evict() // Evict least recently used element if cache is full
	}
	// Add new element to the front of the list
	expireAt := time.Time{}
	if ttl > 0 {
		expireAt = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	newNode := &CacheNode{key: key, value: value, expireAt: expireAt}
	entry := c.list.PushFront(newNode)
	c.cache[key] = entry
}

// evict removes the least recently used element from the cache.
func (c *LRUCache) evict() {
	if evicted := c.list.Back(); evicted != nil {
		c.list.Remove(evicted)
		delete(c.cache, evicted.Value.(*CacheNode).key)
	}
}

// Print returns a string representation of the cache contents in order from most to least recently used.
func (c *LRUCache) Print() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	orderedItems := []string{}
	now := time.Now()
	for elem := c.list.Front(); elem != nil; {
		next := elem.Next()
		node := elem.Value.(*CacheNode)
		if node.expireAt.IsZero() || node.expireAt.After(now) {
			orderedItems = append(orderedItems, fmt.Sprintf("%s:%s", node.key, node.value))
		} else {
			c.list.Remove(elem)
			delete(c.cache, node.key)
		}
		elem = next
	}
	return fmt.Sprintf(strings.Join(orderedItems, ", "))
}

// DEL_ALL deletes the entire cache.
func (c *LRUCache) DEL_ALL() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list.Init()                            // Clear the linked list
	c.cache = make(map[string]*list.Element) // Reset the cache map
}

// Del deletes a key-value pair from the cache.
func (c *LRUCache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.cache[key]; found {
		c.list.Remove(elem)                          // Remove element from linked list
		delete(c.cache, elem.Value.(*CacheNode).key) // Delete from cache map
	}
}
