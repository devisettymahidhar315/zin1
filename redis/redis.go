package redis

import (
	"context"
	"fmt"
	"log"

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
func (c *LRUCache) Put(key, value string, maxLength int) {
	// Check if the key already exists
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Error checking if key %s exists: %v", key, err)
	}
	if exists > 0 {
		// Remove the key from the list if it exists to update its position
		c.client.LRem(ctx, "cache", 0, key)
	}

	// Add the key to the front of the list
	c.client.LPush(ctx, "cache", key)

	// Set the key-value pair in Redis
	c.client.Set(ctx, key, value, 0)

	// Check the current length of the cache
	length, err := c.client.LLen(ctx, "cache").Result()
	if err != nil {
		log.Fatalf("Error getting cache length: %v", err)
	}

	// If the cache exceeds the maxLength, remove the oldest item
	if length > int64(maxLength) {
		oldest, err := c.client.RPop(ctx, "cache").Result()
		if err != nil {
			log.Fatalf("Error popping oldest key: %v", err)
		}
		c.client.Del(ctx, oldest)
	}
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
		if err != nil {
			log.Fatalf("Error getting key %s: %v", key, err)
		}
		orderedItems = append(orderedItems, fmt.Sprintf("%s:%s", key, value))
	}

	// Concatenate the ordered items into a single string
	var result string
	for i, item := range orderedItems {
		result += item
		if i < len(orderedItems)-1 {
			result += ", "
		}
	}

	return result
}

// Del deletes the key-value pair associated with the given key from the cache
func (c *LRUCache) Del(key string) {
	// Check if the key exists
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Error checking if key %s exists: %v", key, err)
	}

	if exists == 0 {
		// Key does not exist
		return
	}

	// Remove the key from the cache list
	_, err = c.client.LRem(ctx, "cache", 0, key).Result()
	if err != nil {
		log.Fatalf("Error removing key %s from cache: %v", key, err)
	}

	// Delete the key-value pair from Redis
	_, err = c.client.Del(ctx, key).Result()
	if err != nil {
		log.Fatalf("Error deleting key %s: %v", key, err)
	}
}

func (c *LRUCache) DEL_ALL() {
	c.client.FlushAll(ctx)
}
