package main

import (
	"sync"
)

type workRequest struct {
	actor string
	path  []relationship
}

type queue struct {
	works []workRequest
	sync.Mutex
}

func (q *queue) enqueue(w ...workRequest) {
	q.Lock()
	defer q.Unlock()
	q.works = append(q.works, w...)
}


func (q *queue) dequeue() workRequest {
	q.Lock()
	defer q.Unlock()
	rel := q.works[0]
	q.works = q.works[1:]
	return rel
}

func (q *queue) empty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.works) == 0
}

func NewQueue() *queue {
	return &queue{}
}