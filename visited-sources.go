package main

import "sync"

// VisitedSources is a set of visited sources with a mutex lock for concurrent access.
// This set is used to avoid repeatedly fetching data from the same source and processing it multiple times.
type VisitedSources struct {
	URLs  map[string]struct{}
	mutex *sync.RWMutex
}

// NewVisitedSources creates a new VisitedSources.
func NewVisitedSources() *VisitedSources {
	return &VisitedSources{
		URLs:  make(map[string]struct{}),
		mutex: &sync.RWMutex{},
	}
}

// Add adds a URL to the VisitedSources.
func (vs *VisitedSources) Add(url string) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()
	vs.URLs[url] = struct{}{}
}

// Has checks if a URL is in the VisitedSources.
func (vs *VisitedSources) Has(url string) bool {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	_, ok := vs.URLs[url]
	return ok
}
