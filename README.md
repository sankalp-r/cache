## Cache with eviction policy

**Problem Statement**: Create an in-memory cache in Go for storing String values based on a String key. Expect that the cache will continually have previously unseen keys added to it, and that some keys will be fetched more frequently than others.

## Design details
This cache implementation has following signature:

```
type Cache interface {
// Inserts the provided key/value pair into the cache, making it
// available for future Get() calls.
Put(key, value string)
	
// Returns a value previously provided via Put(). An empty String
// may be returned if the requested data was never inserted or is
// no longer available.
Get(key string) string
}
```
 This cache implementation is based on least-recently-used eviction policy. Normally to store key-value pair, map data-structure is sufficient. 
 But to evict and find least-recently-used key efficiently, this cache implementation relies on map along with doubly-linked-list .
 
 Here how this implementation works:
 * To add a new key-value pair into the cache, a new-node is added at the head of the linked-list. Also, if a key is fetched recently, it will be removed from its current position and moved to the head of the list.
   So all the recently-used keys will always be at the head of the list.
 * Least-recently-used keys will be on the node which are near to the tail of the list.
 * In-memory map holds (key,Node) pairs which helps to find a node in linked-list in constant(O(1)) time. Similarly, adding/removing a node is also achieved in constant(O(1)) time in doubly-linked-list.
 * After the cache's size exceed the capacity, least-recently-used key is removed from the tail of the linked-list to support eviction policy.
 * To ensure thread-safety, `sync.Mutex` is used for performing operations on the cache.


## How to run
* To download and sync the dependencies run the command make deps-update from Makefile.
* To run unit-tests, run make test from Makefile.
* This implementation is meant to be consumed as library(package).
  Given that you have correct access rights on the repo, it can be consumed by running `go get <module_name>` or `go install <module_name>`.

## Limitations and enhancement-scope
* As this is an in-memory implementation, it will not scale well as the data grows large. To solve this issue, persistence like redis can be used by creating another implementation of `Cache` interface.
* Different types of eviction policies can be easily added by implementing the `Cache` interface.
