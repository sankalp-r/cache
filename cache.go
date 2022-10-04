package cache

import (
	"container/list"
	"errors"
	"sync"
)

type Cache interface {
	Put(key, value string)
	Get(key string) string
}

// LRUCache is the least-recently-used cache implementation.
type LRUCache struct {
	// capacity of the cache.
	capacity int
	// cache to hold linked-list node.
	cache map[string]*list.Element
	// list to store key-value pair.
	list *list.List
	sync.Mutex
}

// pair stores key-value pair.
type pair struct {
	Key   string
	Value string
}

func NewCache(capacity int) (Cache, error) {
	return newLRUCache(capacity)
}

func newLRUCache(capacity int) (Cache, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid capacity: capacity should be greater than zero")
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}, nil
}

func (c *LRUCache) Put(key, value string) {
	c.Lock()
	defer c.Unlock()

	// if value exists, move it to front of the list.
	if node, exist := c.cache[key]; exist {
		node.Value = pair{Key: key, Value: value}
		c.list.MoveToFront(node)
	} else {
		// insert new key at the front of the list.
		node = c.list.PushFront(pair{Key: key, Value: value})
		c.cache[key] = node

		// evict key form tail of the list if length exceeds limit.
		if c.list.Len() > c.capacity {
			tail := c.list.Back()
			c.list.Remove(tail)
			deletedPair := tail.Value.(pair)
			delete(c.cache, deletedPair.Key)
		}
	}
}

func (c *LRUCache) Get(key string) string {
	c.Lock()
	defer c.Unlock()
	// if value exist, move it to front of the list.
	if node, exist := c.cache[key]; exist {
		c.list.MoveToFront(node)
		fetchedPair := node.Value.(pair)
		return fetchedPair.Value
	}
	return ""
}
