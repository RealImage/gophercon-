package main

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type Meta struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Details struct {
	Meta
	Role string `json:"role"`
}

type Person struct {
	Meta
	Type   string    `json:"type"`
	Movies []Details `json:"movies"`
}

type Movie struct {
	Meta
	Type string    `json:"type"`
	Cast []Details `json:"cast"`
	Crew []Details `json:"crew"`
}

type Entity interface {
	Person | Movie
}

type QueueData struct {
	URL        string
	Movie      string
	Role       string
	ParentURL  string
	ParentRole string
	Distance   int
}

type HTTPClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

type SyncManager struct {
	wg  *sync.WaitGroup
	mu  *sync.RWMutex
	dos chan int
}

func NewSyncManager() SyncManager {
	return SyncManager{
		wg:  &sync.WaitGroup{},
		mu:  &sync.RWMutex{},
		dos: make(chan int), // degrees of separation
	}
}

type Path struct {
	ParentURL  string // Parent URL
	Movie      string
	Role       string // Role of the Main Actor
	ParentRole string // Role of the parent actor
}
