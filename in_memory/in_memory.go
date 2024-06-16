package in_memory

import (
	"container/list"
	"fmt"
)

// CacheNode represents a single node in the LRU cache with a key-value pair.
type CacheNode struct {
	key   string
	value string
}

// LRUCache implements a Least Recently Used (LRU) cache using a map and a doubly linked list.
type LRUCache struct {
	cache map[string]*list.Element
	list  *list.List
}

// NewLRUCache initializes and returns a new LRUCache instance.
func NewLRUCache() *LRUCache {
	return &LRUCache{
		cache: make(map[string]*list.Element),
		list:  list.New(),
	}
}

// Get retrieves the value associated with the given key.
// It moves the accessed element to the front of the list to mark it as recently used.
func (c *LRUCache) Get(key string) string {
	if elem, found := c.cache[key]; found {
		c.list.MoveToFront(elem)
		return elem.Value.(*CacheNode).value
	}
	return ""
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves the element to the front.
// If the cache exceeds maxLength, it evicts the least recently used element.
func (c *LRUCache) Put(key string, value string, maxLength int) {
	if elem, found := c.cache[key]; found {
		elem.Value.(*CacheNode).value = value
		c.list.MoveToFront(elem)
		return
	}
	if c.list.Len() == maxLength {
		// Evict the least recently used element
		evicted := c.list.Back()
		if evicted != nil {
			c.list.Remove(evicted)
			delete(c.cache, evicted.Value.(*CacheNode).key)
		}
	}
	// Add the new element to the front of the list
	newNode := &CacheNode{key: key, value: value}
	entry := c.list.PushFront(newNode)
	c.cache[key] = entry
}

// Del deletes the key-value pair associated with the given key from the cache.
func (c *LRUCache) Del(key string) {
	if elem, found := c.cache[key]; found {
		c.list.Remove(elem)
		delete(c.cache, elem.Value.(*CacheNode).key)
	}
}

// Print returns a string representation of the cache contents in order from most to least recently used.
func (c *LRUCache) Print() string {
	orderedItems := []string{}

	// Iterate through the list from front to back
	for elem := c.list.Front(); elem != nil; elem = elem.Next() {
		node := elem.Value.(*CacheNode)
		orderedItems = append(orderedItems, fmt.Sprintf("%s:%s", node.key, node.value))
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

func (c *LRUCache) DEL_ALL() {
	for elem := c.list.Front(); elem != nil; {
		next := elem.Next() // Store the next element before removing the current one
		c.list.Remove(elem)
		delete(c.cache, elem.Value.(*CacheNode).key)
		elem = next // Move to the next element
	}
}
