package main

import "sync"

type VisitedSources struct {
	URLs  map[string]struct{}
	mutex *sync.RWMutex
}

func NewVisitedSources() VisitedSources {
	return VisitedSources{
		URLs:  make(map[string]struct{}),
		mutex: &sync.RWMutex{},
	}
}

func (vs *VisitedSources) Add(url string) {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()
	vs.URLs[url] = struct{}{}
}

func (vs *VisitedSources) Has(url string) bool {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	_, ok := vs.URLs[url]
	return ok
}

var visitedSources = NewVisitedSources()
