package queue

import (
	"container/ring"
	"errors"
)

// List - a Queue backed by a container/list
type RingQueue struct {
	ring *ring.Ring	 // preallocated ring for all slots in the queue
	head *ring.Ring  // head of the queue
	tail *ring.Ring  // tail of the queue
	capacity int     // maximum number of elements the ring can hold
	length int       // current number of element in the ring
}

func (rq *RingQueue) Len() int {
	return rq.length
}

func (rq *RingQueue) Cap() int {
	return rq.capacity
}

func (rq *RingQueue) Push(value interface{}) error {
	if rq.length >= rq.capacity {
		return errors.New("queue is full")
	}
	// insert at tail
	rq.tail.Value = value

	// increment the tail
	rq.tail = rq.tail.Next()

	// increment the length
	rq.length++

	return nil
}

func (rq *RingQueue) Pop() (interface{}, error) {
	if rq.length == 0 {
		return nil, errors.New("queue is empty)")
	}

	// get the value at the head
	value := rq.head.Value

	// increment the head
	rq.head = rq.head.Next()

	// decrement length
	rq.length--

	return value,nil
}

func NewRingQueue(cap int) Queue {
	var rq RingQueue
	rq.capacity = cap
	rq.length = 0
	rq.ring = ring.New(cap)
	rq.head = rq.ring
	rq.tail = rq.ring

	return &rq
}

func NewSyncRing(cap int) BoundedQueue {
	var rq Queue
	var bq BoundedQueue 

	rq = NewRingQueue(cap)

	bq = NewQueueSync(rq) 

	return bq
}