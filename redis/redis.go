package redis

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// Create a background context for Redis operations
var ctx = context.Background()

// LRUCache represents a Redis-based LRU cache
type LRUCache struct {
	client *redis.Client
}

// NewLRUCache initializes and returns a new LRUCache instance connected to Redis
func NewLRUCache() *LRUCache {
	opts := &redis.Options{
		Addr:     "localhost:6379", // Redis server address
		PoolSize: 10,               // Connection pool size
	}
	// Create a new Redis client
	rdb := redis.NewClient(opts)
	// Clear the cache on initialization
	rdb.Del(ctx, "cache")

	return &LRUCache{
		client: rdb,
	}
}

// Put adds or updates a key-value pair in the cache
// If the cache exceeds maxLength, the least recently used item is removed
func (c *LRUCache) Put(key, value string, maxLength, ttl int) {
	// Check if the key already exists
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Error checking if key %s exists: %v", key, err)
	}
	if exists > 0 {
		// Remove the key from the list to update its position
		c.client.LRem(ctx, "cache", 0, key)
	}
	// Add the key to the front of the list
	c.client.LPush(ctx, "cache", key)

	if ttl == -1 {
		// No expiration
		c.client.Set(ctx, key, value, 0)
	} else if ttl > 0 {
		// Set key with expiration time
		c.client.Set(ctx, key, value, time.Duration(ttl)*time.Second)
	} else {
		log.Fatalf("Invalid TTL value: %d. TTL should be -1 (no expiration) or greater than 0", ttl)
	}
	// Ensure cache size does not exceed maxLength
	c.evictItems(maxLength)
}

// Get retrieves the value associated with the given key
// If the key is found, it is moved to the front of the list
func (c *LRUCache) Get(key string) string {
	// Get the value associated with the key
	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "" // Key does not exist
	} else if err != nil {
		log.Fatalf("Error getting key %s: %v", key, err)
	}
	// Move the key to the front of the list
	c.client.LRem(ctx, "cache", 0, key)
	c.client.LPush(ctx, "cache", key)

	return value
}

// Print returns a string representation of the cache contents
func (c *LRUCache) Print() string {
	// Get all keys from the cache list
	keys, err := c.client.LRange(ctx, "cache", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting cache keys: %v", err)
	}
	orderedItems := []string{}

	// Retrieve the values for each key and format them
	for _, key := range keys {
		value, err := c.client.Get(ctx, key).Result()
		if err == redis.Nil {
			// Key does not exist, remove it from the list
			c.client.LRem(ctx, "cache", 0, key)
			continue
		} else if err != nil {
			log.Fatalf("Error getting key %s: %v", key, err)
		}
		orderedItems = append(orderedItems, fmt.Sprintf("%s:%s", key, value))
	}
	// Concatenate the ordered items into a single string
	return strings.Join(orderedItems, ", ")
}

// Del deletes the key-value pair associated with the given key from the cache
func (c *LRUCache) Del(key string) {
	// Check if the key exists
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Error checking if key %s exists: %v", key, err)
	}

	if exists > 0 {
		// Remove the key from the cache list
		if _, err := c.client.LRem(ctx, "cache", 0, key).Result(); err != nil {
			log.Fatalf("Error removing key %s from cache: %v", key, err)
		}
		// Delete the key-value pair from Redis
		if _, err := c.client.Del(ctx, key).Result(); err != nil {
			log.Fatalf("Error deleting key %s: %v", key, err)
		}
	}
}

// DEL_ALL clears the entire cache
func (c *LRUCache) DEL_ALL() {
	c.client.FlushAll(ctx)
}

// evictItems ensures the cache size does not exceed maxLength
func (c *LRUCache) evictItems(maxLength int) {
	// Get current length of the cache
	length, err := c.client.LLen(ctx, "cache").Result()
	if err != nil {
		log.Fatalf("Error getting cache length: %v", err)
	}

	// Check for expired keys and remove them
	for i := int64(0); i < length; i++ {
		keyToCheck, err := c.client.LIndex(ctx, "cache", i).Result()
		if err != nil {
			log.Fatalf("Error getting key at index %d: %v", i, err)
		}
		exists, err := c.client.Exists(ctx, keyToCheck).Result()
		if err != nil {
			log.Fatalf("Error checking if key %s exists: %v", keyToCheck, err)
		}
		if exists == 0 {
			c.client.LRem(ctx, "cache", 0, keyToCheck)
			i--
			length--
		}
	}

	// Evict excess items if cache size exceeds maxLength
	for length > int64(maxLength) {
		oldest, err := c.client.RPop(ctx, "cache").Result()
		if err != nil {
			log.Fatalf("Error popping oldest key: %v", err)
		}
		c.client.Del(ctx, oldest)
		length, err = c.client.LLen(ctx, "cache").Result()
		if err != nil {
			log.Fatalf("Error getting cache length: %v", err)
		}
	}
}
