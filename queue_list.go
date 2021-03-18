package queue

import (
	"container/list"
	"errors"
)

// ListQueue - a Queue backed by a container/list
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

func NewSyncList(cap int) BoundedQueue {
	var lq Queue
	var bq BoundedQueue 

	lq = NewListQueue(cap)

	bq = NewQueueSync(lq) 

	return bq
}
