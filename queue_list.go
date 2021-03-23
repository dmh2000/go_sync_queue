package queue

import (
	"container/list"
	"errors"
)

// ListQueue backed by container/list
// the container/list data structure supports the semantics and methods
// needed for the Queue interface, with the exception of a capacity.
type ListQueue struct {
	list *list.List	 // contains the elements currently in the queue
	capacity int	     // maximum number of elements the queue can hold
}

func (lq *ListQueue) Len() int {
	return lq.list.Len()
}

func (lq *ListQueue) Cap() int {
	return lq.capacity
}

func (lq *ListQueue) Push(value interface{}) error {
	if lq.list.Len() >= lq.capacity {
		return errors.New("queue is full")
	}
	// insert 
	lq.list.PushBack(value)

	return nil
}

func (lq *ListQueue) Pop() (interface{}, error) {
	if lq.list.Len() == 0 {
		return nil, errors.New("queue is empty)")
	}

	value := lq.list.Remove(lq.list.Front())

	return value,nil
}

func NewListQueue(cap int) Queue {
	var lq ListQueue
	lq.capacity = cap
	lq.list = list.New()

	return &lq
}

func NewSyncList(cap int) SynchronizedQueue {
	var lq Queue
	var bq SynchronizedQueue 

	lq = NewListQueue(cap)

	bq = NewSynchronizedQueue(lq) 

	return bq
}
