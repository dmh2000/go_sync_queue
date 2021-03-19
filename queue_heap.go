package queue

import (
	"container/heap"
	"errors"
)

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// see priority_queue.go in this directory for the definition of the Item struct
// that are the members of this queue
// HeapQueue - a Queue backed by a container/heap - PriorityQueue example
type HeapQueue struct {
	heap IntHeap	 // use the priority queue example in priority_queue.go
	capacity int	     // maximum number of elements the queue can hold
}

func (pq HeapQueue) Len() int {
	return pq.heap.Len()
}

func (pq HeapQueue) Cap() int {
	return pq.capacity
}

func (pq *HeapQueue) Push(value interface{}) error {
	if pq.heap.Len() >= pq.capacity {
		return errors.New("queue is full")
	}
	// insert 
	heap.Push(&pq.heap, value)

	return nil
}

func (pq *HeapQueue) Pop() (interface{}, error) {
	if pq.heap.Len() == 0 {
		return nil, errors.New("queue is empty)")
	}

	value := heap.Pop(&pq.heap)

	return value,nil
}

// create a new heap queue
func NewIntHeap(cap int) Queue {
	var pq HeapQueue

	// set the capacity
	pq.capacity = cap

	// set up an empty heap to start with
	pq.heap = make(IntHeap,0)

	// initialize it
	heap.Init(&pq.heap)

	return &pq
}

// wrap the heap queue in a Synchronized queue
func NewSyncHeap(cap int) SynchronizedQueue {
	var pq Queue
	var bq SynchronizedQueue 

	// create the int heap
	pq = NewIntHeap(cap)

	// wrap it in the syncrhonized bounded queue
	bq = NewSynchronizedQueue(pq) 

	return bq
}
