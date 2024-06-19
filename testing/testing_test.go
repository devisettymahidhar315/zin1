package testing

import (
	"testing"

	"github.com/devisettymahidhar315/zin1/in_memory"
	"github.com/devisettymahidhar315/zin1/multi_cache"
	"github.com/devisettymahidhar315/zin1/redis"
)

// Length of the cache for testing
// in memory
const len1 = 2

// TestPut_inmemory tests the Put method of the inmemory cache
func TestPut_inmemory(t *testing.T) {
	// Initialize a new inmemory cache
	cache := in_memory.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Retrieve the value for key "a1"
	result := cache.Get("a1")

	// Check if the value is as expected
	if result != "1" {
		t.Error("Expected value '1', got", result)
	}
}

// TestGet_inmemory tests the Get method of the inmemory cache

func TestGet_inmemory(t *testing.T) {
	// Initialize a new inmemory cache
	cache := in_memory.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Retrieve the value for key "a1"
	result := cache.Get("a1")
	if result != "1" {
		t.Error("Expected value '1', got", result)
	}

	// Attempt to retrieve a value for a non-existent key "v"
	result = cache.Get("v")
	if result != "" {
		t.Error("Expected empty string for non-existent key, got", result)
	}
}

// TestPrint_inmemory tests the Print method of the inmemory cache
func TestPrint_inmemory(t *testing.T) {
	// Initialize a new inmemory cache
	cache := in_memory.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Print the current state of the cache
	result := cache.Print()
	expected_result := "b1:2, a1:1"

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}

	// Insert another key-value pair to exceed the cache length
	cache.Put("c1", "3", len1)

	// Print the current state of the cache
	result = cache.Print()
	expected_result = "c1:3, b1:2"

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}
}

// TestDel_inmemory tests the Del method of the inmemory cache
func TestDel_inmemory(t *testing.T) {
	// Initialize a new inmemory cache
	cache := in_memory.NewLRUCache()

	// Insert a key-value pair into the cache
	cache.Put("a1", "1", len1)

	// Delete the key "a1"
	cache.Del("a1")

	// Print the current state of the cache
	result := cache.Print()
	expected_result := ""

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Delete the key "a1"
	cache.Del("a1")

	// Attempt to retrieve the value for the deleted key "a1"
	result = cache.Get("a1")
	expected_result = ""

	// Check if the result matches the expected result
	if result != expected_result {
		t.Error("Expected empty string for deleted key, got", result)
	}
}

// TestDelAll_inmemory tests the Del_All method of the inmemory cache

func TestDelAll_inmemory(t *testing.T) {
	// Initialize a new inmemory  cache
	cache := in_memory.NewLRUCache()
	// Insert a key-value pair into the cache

	cache.Put("a", "1", len1)
	cache.Put("b", "2", len1)

	// Delete all key
	cache.DEL_ALL()

	// Attempt to retrieve the value for the deleted key "a1"
	result := cache.Get("a1")
	expected_result := ""

	// Check if the result matches the expected result
	if result != expected_result {
		t.Error("Expected empty string for deleted key, got", result)
	}

}

// redis
// TestPut_redis tests the Put method of the Redis cache
func TestPut_redis(t *testing.T) {
	// Initialize a new Redis cache
	cache := redis.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Retrieve the value for key "a1"
	result := cache.Get("a1")

	// Check if the value is as expected
	if result != "1" {
		t.Error("Expected value '1', got", result)
	}
}

// TestGet_redis tests the Get method of the Redis cache
func TestGet_redis(t *testing.T) {
	// Initialize a new Redis cache
	cache := redis.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Retrieve the value for key "a1"
	result := cache.Get("a1")
	if result != "1" {
		t.Error("Expected value '1', got", result)
	}

	// Attempt to retrieve a value for a non-existent key "v"
	result = cache.Get("v")
	if result != "" {
		t.Error("Expected empty string for non-existent key, got", result)
	}
}

// TestPrint_redis tests the Print method of the Redis cache
func TestPrint_redis(t *testing.T) {
	// Initialize a new Redis cache
	cache := redis.NewLRUCache()

	// Insert key-value pairs into the cache
	cache.Put("a1", "1", len1)
	cache.Put("b1", "2", len1)

	// Print the current state of the cache
	result := cache.Print()

	expected_result := "b1:2, a1:1"

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}

	// Insert another key-value pair to exceed the cache length
	cache.Put("c1", "3", len1)

	// Print the current state of the cache
	result = cache.Print()
	expected_result = "c1:3, b1:2"

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}
}

// TestDel_redis tests the Del method of the Redis cache
func TestDel_redis(t *testing.T) {
	// Initialize a new Redis cache
	cache := redis.NewLRUCache()

	// Insert a key-value pair into the cache
	cache.Put("a1", "1", len1)

	// Delete the key "a1"
	cache.Del("a1")

	// Print the current state of the cache
	result := cache.Print()
	expected_result := ""

	// Check if the printed result matches the expected result
	if result != expected_result {
		t.Error("Expected", expected_result, "got", result)
	}

}

// TestDelAll_redis tests the Del_All method of the Redis cache
func TestDelAll_redis(t *testing.T) {
	// Initialize a new redis  cache
	cache := redis.NewLRUCache()
	// Insert a key-value pair into the cache

	cache.Put("a", "1", len1)
	cache.Put("b", "2", len1)

	// Delete all key
	cache.DEL_ALL()

	// Attempt to retrieve the value for the deleted key "a1"
	result := cache.Get("a1")
	expected_result := ""

	// Check if the result matches the expected result
	if result != expected_result {
		t.Error("Expected empty string for deleted key, got", result)
	}

}

// integration testing
// testing related to multi_cache folder
func TestGET(t *testing.T) {
	// Create a new multi-cache instance
	cache := multi_cache.NewMultiCache()
	// Set key-value pairs in the cache
	cache.Set("a", "1", len1)
	cache.Set("b", "2", len1)

	// Test case 1: Get a non-existent key
	res1 := cache.Get("c")
	if res1 != "" {
		t.Error("case 1 error: expected empty string for non-existent key")
	}

	// Test case 2: Get an existing key
	res2 := cache.Get("a")
	if res2 != "1" {
		t.Error("case 2 error: expected '1' for key 'a'")
	}
}

// TestPrint tests the Print methods of the multi_cache
func TestPrint(t *testing.T) {
	// Create a new multi-cache instance
	cache := multi_cache.NewMultiCache()
	// Set key-value pairs in the cache
	cache.Set("a", "1", len1)
	cache.Set("b", "2", len1)

	// Get the printed results from both in-memory and Redis caches
	inmemory_result := cache.Print_in_mem()
	redis_result := cache.Print_redis()
	// Check if the data is the same in both backends
	if inmemory_result != redis_result {
		t.Error("data is not the same in both backends")
	}

	// Set another key-value pair to exceed the cache capacity
	cache.Set("c", "3", len1)
	inmemory_result = cache.Print_in_mem()
	redis_result = cache.Print_redis()

	// Check if the data is the same in both backends
	if inmemory_result != redis_result {
		t.Error("data is not the same in both backends")
	}

}

// TestDel tests the Del method of the cache
func TestDel(t *testing.T) {
	// Create a new multi-cache instance
	cache := multi_cache.NewMultiCache()
	// Set key-value pairs in the cache
	cache.Set("a", "1", len1)
	cache.Set("b", "2", len1)

	// Delete a key from the cache
	cache.Del("a")
	// Check if the deleted key returns an empty string
	result := cache.Get("a")
	if result != "" {
		t.Error("expected empty string for deleted key 'a'")
	}
	// Check if an existing key returns the correct value
	result = cache.Get("b")
	if result != "2" {
		t.Error("expected '2' for key 'b'")
	}
}
