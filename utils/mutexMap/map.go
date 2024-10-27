package mutexmap

import "sync"

type CallVisitedMutex struct {
	sync.Mutex
	visited map[string]bool
}

func NewCallVisitedMutex() *CallVisitedMutex {
	c := CallVisitedMutex{
		visited: make(map[string]bool),
	}

	return &c
}

func (c *CallVisitedMutex) Set(name string) {
	c.Lock()
	c.visited[name] = true
	c.Unlock()
}

func (c *CallVisitedMutex) Get(name string) bool {
	c.Lock()
	set := c.visited[name]
	c.Unlock()

	return set
}
