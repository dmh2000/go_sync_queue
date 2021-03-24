package queue

import (
	"container/heap"
	"errors"
	"fmt"
)

// PriorityItem - for a type agnostic priority queue
// modeled after the PriorityQueue Example at
// https://golang.org/pkg/container/heap/#example__priorityQueue
type PriorityItem struct {
	value interface{}
	priority int
}

type PrioritySlice []PriorityItem

func (h PrioritySlice) Len() int           { return len(h) }
func (h PrioritySlice) Less(i, j int) bool { return h[i].priority < h[j].priority }
func (h PrioritySlice) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PrioritySlice) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(PriorityItem))
}

func (h *PrioritySlice) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PriorityQueue - a Queue backed by a container/heap - PriorityQueue example
type PriorityQueue struct {
	heap PrioritySlice	 // use the priority queue example in priority_queue.go
	capacity int	     // maximum number of elements the queue can hold
}

func (pq PriorityQueue) Len() int {
	return pq.heap.Len()
}

func (pq PriorityQueue) Cap() int {
	return pq.capacity
}

func (pq *PriorityQueue) Push(value interface{}) error {
	if pq.heap.Len() >= pq.capacity {
		return errors.New("queue is full")
	}
	// insert 
	heap.Push(&pq.heap, value)

	return nil
}

func (pq *PriorityQueue) Pop() (interface{}, error) {
	if pq.heap.Len() == 0 {
		return nil, errors.New("queue is empty)")
	}

	value := heap.Pop(&pq.heap)

	return value,nil
}

// String
func (pq *PriorityQueue) String() string {
	return fmt.Sprintf("PriorityQueue Len:%v Cap:%v",pq.Len(),pq.Cap())
}

// create a new heap queue
func NewPriorityQueue(cap int) Queue {
	var pq PriorityQueue

	// set the capacity
	pq.capacity = cap

	// set up an empty heap to start with
	pq.heap = make(PrioritySlice,0)

	// initialize it
	heap.Init(&pq.heap)

	return &pq
}

// wrap the heap queue in a Synchronized queue
func NewSyncPriority(cap int) SynchronizedQueue {
	var pq Queue
	var bq SynchronizedQueue 

	// create the int heap
	pq = NewPriorityQueue(cap)

	// wrap it in the syncrhonized bounded queue
	bq = NewSynchronizedQueue(pq) 

	return bq
}
