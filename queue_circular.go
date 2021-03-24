package queue

import (
	"errors"
	"fmt"
)

// Implementation of Queue interface using circular buffer
type CircularQueue struct {
	queue []interface{}	// data
	head     int		// items are pulled from the head
	tail     int		// items are pushed to the tail
	length   int		// current number of elements in the queue
	capacity int	    // maximum allowed elements total
}

func (cb *CircularQueue) Len() int {
	return cb.length
}

func (cb *CircularQueue) Cap() int {
	return cb.capacity
}

func (cb *CircularQueue) Push(value interface{}) error {
	if cb.length >= cb.capacity {
		return errors.New("queue is full")
	}
	// insert and count
	cb.queue[cb.tail] = value
	cb.tail = (cb.tail+1) % cb.capacity
	cb.length++

	return nil
}

func (cb *CircularQueue) Pop() (interface{}, error) {
	if cb.length == 0 {
		return nil, errors.New("queue is empty)")
	}
	value := cb.queue[cb.head]
	cb.head = (cb.head + 1)  % cb.capacity
	cb.length--

	return value,nil
}

// String
func (cb *CircularQueue) String() string {
	return fmt.Sprintf("CircularQueue Len:%v Cap:%v",cb.Len(),cb.Cap())
}


func NewCircularQueue(cap int) Queue {
	var cq CircularQueue

	cq.length = 0
	cq.capacity = cap
	cq.head = 0
	cq.tail = 0
	cq.queue = make([]interface{},cap)

	return &cq
}

func NewSyncCircular(cap int) SynchronizedQueue {
	var cq Queue
	var bq SynchronizedQueue 

	cq = NewCircularQueue(cap)

	bq = NewSynchronizedQueue(cq) 

	return bq
}
