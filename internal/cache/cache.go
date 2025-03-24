// Package cache provides an in-memory storage system for caching data in the application.
// It includes operations for loading, reading, and checking data in the cache, using synchronization mechanisms
// to ensure thread-safe access to the cached data.
package cache

import (
	"sber/pkg/models"
	"sync"
	"sync/atomic"
)

// Storage represents the in-memory cache storage system. It contains a map for storing cached data
// and an atomic counter for generating unique IDs for each cache entry. The struct also uses a mutex to ensure
// thread-safety when accessing and modifying the cache.
type Storage struct {
	// str is a map that stores the cached data with an int64 key (ID) and a CacheStorageFormat value.
	str map[int32]models.CacheStorageFormat
	// mu is a Mutex used to synchronize access to the cache.
	mu sync.Mutex
	// IDCounter is an atomic counter used to generate unique IDs for cache entries.
	IDCounter int32
}

// New creates and returns a new instance of the Storage struct with an empty cache map.
func New() *Storage {
	return &Storage{
		str: map[int32]models.CacheStorageFormat{},
	}
}

// Load adds a new entry to the cache with a unique ID and the given value. It increments the ID counter atomically
// to ensure that each entry gets a unique ID.
func (s *Storage) Load(value models.Result) {
	// Load the current value of IDCounter atomically to generate a unique ID.
	id := atomic.LoadInt32(&s.IDCounter)
	// Increment the IDCounter atomically.
	atomic.AddInt32(&s.IDCounter, 1)

	// Lock the mutex to ensure thread-safe access to the cache while modifying it.
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the new cache entry with the generated ID.
	cacheData := models.CacheStorageFormat{
		ID:         id,
		Params:     value.Params,
		Program:    value.Program,
		Aggregates: value.Aggregates,
	}

	cacheData.MarshalJSON()
	s.str[id] = cacheData
}

// ReadAll returns all entries from the cache as a slice of CacheStorageFormat. It locks the cache before reading
// to ensure thread-safety.
func (s *Storage) ReadAll() []models.CacheStorageFormat {
	// Lock the mutex to ensure thread-safe access to the cache while reading it.
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a slice to hold the cache entries.
	strArr := make([]models.CacheStorageFormat, len(s.str))

	// Copy the cache entries into the slice.
	i := 0
	for _, v := range s.str {
		strArr[i] = v
		i++
	}

	// Return the slice containing all cache entries.
	return strArr
}

// HasData checks whether there are any entries in the cache. It returns true if the cache is not empty.
func (s *Storage) HasData() bool {
	// Return whether the cache map is empty or not.
	return len(s.str) != 0
}
