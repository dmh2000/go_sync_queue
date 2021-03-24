package queue

import (
	"errors"
	"fmt"
)

// SliceQueue backed by a slice
// this version appends and removes elements so the slice grows and shrinks
// memory is not preallocated
// see CircularQueue for a version that preallocates a slice of capacity
type SliceQueue struct {
	slice []interface{}
	capacity int	     // maximum number of elements the queue can hold
}

func (sq *SliceQueue) Len() int {
	return len(sq.slice)
}

func (sq *SliceQueue) Cap() int {
	return sq.capacity
}

func (sq *SliceQueue) Push(value interface{}) error {
	if len(sq.slice) >= sq.capacity {
		return errors.New("queue is full")
	}
	// insert at end
	sq.slice = append(sq.slice, value)

	return nil
}

func (sq *SliceQueue) Pop() (interface{}, error) {
	if len(sq.slice) == 0 {
		return nil, errors.New("queue is empty)")
	}

	// get the value at the front
	value := sq.slice[0]

	// remove the front
	sq.slice = sq.slice[1:]


	return value,nil
}

// String
func (sq *SliceQueue)  String() string {
	return fmt.Sprintf("SliceQueue Len:%v Cap:%v",sq.Len(),sq.Cap())
}


func NewSliceQueue(cap int) Queue {
	var sq SliceQueue
	sq.capacity = cap
	sq.slice = make([]interface{},0)

	return &sq
}

func NewSyncSlice(cap int) SynchronizedQueue {
	var sq Queue
	var bq SynchronizedQueue 

	sq = NewSliceQueue(cap)

	bq = NewSynchronizedQueue(sq) 

	return bq
}
