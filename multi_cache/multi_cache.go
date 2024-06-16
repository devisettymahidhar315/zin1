package multi_cache

import (
	"sync"

	"github.com/devisettymahidhar315/zin1/in_memory"
	"github.com/devisettymahidhar315/zin1/redis"
)

// MultiCache struct manages both Redis and in-memory caches.
type MultiCache struct {
	redisCache    *redis.LRUCache
	inMemoryCache *in_memory.LRUCache
}

// NewMultiCache initializes a new MultiCache with Redis and in-memory LRU caches.
func NewMultiCache() *MultiCache {
	return &MultiCache{
		redisCache:    redis.NewLRUCache(),
		inMemoryCache: in_memory.NewLRUCache(),
	}
}

// Set stores the key-value pair in both Redis and in-memory caches concurrently.
func (c *MultiCache) Set(key, value string, length int) {
	var wg sync.WaitGroup
	wg.Add(2) // Add two goroutines to the wait group
	// Store in Redis cache concurrently
	go func() {
		defer wg.Done()
		c.redisCache.Put(key, value, length)
	}()
	// Store in inmemory cache concurrently
	go func() {
		defer wg.Done()
		c.inMemoryCache.Put(key, value, length)
	}()
	wg.Wait() // Wait for both goroutines to finish
}

// Get retrieves the value for a key from both Redis and in-memory caches concurrently and compares them.
func (c *MultiCache) Get(key string) string {
	var redis_value, inmemory_value string
	var wg sync.WaitGroup
	wg.Add(2) // Add two goroutines to the wait group
	// Retrieve value from in-memory cache concurrently
	go func() {
		defer wg.Done()
		inmemory_value = c.inMemoryCache.Get(key)
	}()
	// Retrieve value from redis cache concurrently
	go func() {
		defer wg.Done()
		redis_value = c.redisCache.Get(key)
	}()
	wg.Wait() // Wait for both goroutines to finish

	// Return the value if they match, otherwise return an empty string
	if redis_value == inmemory_value {
		return redis_value
	} else {
		return ""
	}

}

// Print_redis prints the contents of the Redis cache.
func (c *MultiCache) Print_redis() string {
	var wg sync.WaitGroup
	wg.Add(1) // Add one goroutine to the wait group
	var result string

	// Print the Redis cache contents concurrently
	go func() {
		defer wg.Done()
		result = c.redisCache.Print()
	}()

	wg.Wait() // Wait for the goroutine to finish
	return result
}

// Print_in_mem prints the contents of the in-memory cache.
func (c *MultiCache) Print_in_mem() string {
	var wg sync.WaitGroup
	wg.Add(1) // Add one goroutine to the wait group
	var result string

	// Print the in-memory cache contents concurrently
	go func() {
		defer wg.Done()
		result = c.inMemoryCache.Print()
	}()

	wg.Wait() // Wait for the goroutine to finish
	return result
}

// Del deletes the key-value pair from both Redis and in-memory caches concurrently.
func (c *MultiCache) Del(key string) {
	var wg sync.WaitGroup
	wg.Add(2) // Add two goroutines to the wait group

	// Delete from in-memory cache concurrently
	go func() {
		defer wg.Done()
		c.inMemoryCache.Del(key)
	}()

	// Delete from Redis cache concurrently
	go func() {
		defer wg.Done()
		c.redisCache.Del(key)
	}()

	wg.Wait() // Wait for both goroutines to finish
}

func (c *MultiCache) Del_ALL() {
	var wg sync.WaitGroup
	wg.Add(2) // Add two goroutines to the wait group

	// Delete from in-memory cache concurrently
	go func() {
		defer wg.Done()
		c.inMemoryCache.DEL_ALL()
	}()

	// Delete from Redis cache concurrently
	go func() {
		defer wg.Done()
		c.redisCache.DEL_ALL()
	}()

	wg.Wait() // Wait for both goroutines to finish
}
