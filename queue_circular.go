package queue

import (
	"errors"
)

// Implementation of Queue interface using circular buffer
type CircularBuffer struct {
	queue []interface{}
	head     int
	tail     int
	length   int
	capacity int
}

func (cb *CircularBuffer) Len() int {
	return cb.length
}

func (cb *CircularBuffer) Cap() int {
	return cb.capacity
}

func (cb *CircularBuffer) Push(value interface{}) error {
	if cb.length >= cb.capacity {
		return errors.New("queue is full")
	}
	// insert and count
	cb.queue[cb.tail] = value
	cb.tail = (cb.tail+1) % cb.capacity
	cb.length++

	return nil
}

func (cb *CircularBuffer) Pop() (interface{}, error) {
	if cb.length == 0 {
		return nil, errors.New("queue is empty)")
	}
	value := cb.queue[cb.head]
	cb.head = (cb.head + 1)  % cb.capacity
	cb.length--

	return value,nil
}

func NewCircularBuffer(cap int) Queue {
	var cq CircularBuffer

	cq.length = 0
	cq.capacity = cap
	cq.head = 0
	cq.tail = 0
	cq.queue = make([]interface{},cap)

	return &cq
}

func NewSyncCircular(cap int) BoundedQueue {
	var cq Queue
	var bq BoundedQueue 

	cq = NewCircularBuffer(cap)

	bq = NewQueueSync(cq) 

	return bq
}
