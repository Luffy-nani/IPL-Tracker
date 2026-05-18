package main

import (
	"sync"
	"time"
)

// the fields inside the cache struct are lowercase which means they are private...we dont all the users to access the fields as it can cause race conditions...we use mutexes on the fucntion and users can set and get data from it
type Cache struct {
	matches   []Match
	updatedAt time.Time
	mu        sync.RWMutex
}

var cache = &Cache{}

func (c *Cache) Set(matches []Match) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.matches = matches
	c.updatedAt = time.Now()

}

func (c *Cache) Get() ([]Match, time.Time) {
	c.mu.RLock() // we use Rlock because multiple readers can access the buffer its fine(reader writers problem)
	defer c.mu.RUnlock()
	return c.matches, c.updatedAt //we can return mutliple values in go unlike in cpp
}

func (c *Cache) GetByID(id string) (Match, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, match := range c.matches {
		if match.ID == id {
			return match, true
		}
	}
	return Match{}, false
}
