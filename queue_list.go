package queue

import (
	"container/list"
	"errors"
	"sync"
)

// List is a type of queue that uses a
// condition variable and lists to implement the
// BoundedQueue interface. this implementation
// is intended to be thread safe
type List struct {
	queue *list.List		// contains the elements currently in the queue
	capacity int	    // maximum number of elements the queue can hold
	mtx sync.Mutex      // a mutex for mutual exclusion
	cvr *sync.Cond       // a condition variable for controlling mutations to the queue
}

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (cvq *List) Put(value interface{})  {
	// local the mutex
	cvq.cvr.L.Lock()
	defer cvq.cvr.L.Unlock()


	// block until a value is in the queue
	for cvq.queue.Len() == cvq.capacity {
		// releast and wait
		cvq.cvr.Wait()
	}
	
	// queue had room, add it
	cvq.queue.PushBack(value)

	// signal a waiter if any
	cvq.cvr.Signal()
} 

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (cvq *List) TryPut(value interface{}) error {
	// local the mutex
	cvq.cvr.L.Lock();
	defer cvq.cvr.L.Unlock()

	// is queue full ?
	if cvq.queue.Len() == cvq.capacity {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it
	cvq.queue.PushBack(value)

	// signal a waiter if any
	cvq.cvr.Signal()

	// no error
	return nil
} 



// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (cvq *List) Get() interface{} {
	// lock the mutex
	cvq.cvr.L.Lock()
	defer cvq.cvr.L.Unlock()

	// block until a value is in the queue
	for cvq.queue.Len() == 0 {
		// releast and wait
		cvq.cvr.Wait()
	}

	// at this point there is at least one item in the queue
	// remove the head
	value := cvq.queue.Remove(cvq.queue.Front())

	// unlock the mutex
	return value
}

// TryGet gets a value or returns an error if the queue is empty
func (cvq *List) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	cvq.cvr.L.Lock()
	defer cvq.cvr.L.Unlock()

	// is the queue empty?
	if cvq.queue.Len() > 0 {
		value = cvq.queue.Remove(cvq.queue.Front())
	} else {
		value = nil
		err = errors.New("queue is empty");
	}
	
	return value, err
}

// Len is the current number of elements in the queue 
func (cvq *List) Len() int {
	return cvq.queue.Len()
}

// Cap is the maximum number of elements the queue can hold
func (cvq *List) Cap() int {
	return cvq.capacity;
}

// String
func (cvq *List) String() string {return ""}

// NewListQueue is a factory for creating bounded queues
// that use a condition variable and lists. It returns
// an instance of pointer to BoundedQueue
func NewListQueue(size int) BoundedQueue {
	var cvq List

	cvq.capacity = size
	cvq.queue = list.New()
	cvq.mtx = sync.Mutex{} // unlock mutex
	cvq.cvr = sync.NewCond(&cvq.mtx)

	return &cvq
}

