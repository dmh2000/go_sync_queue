package queue

import (
	"errors"
	"sync"
)

// CondQ2 is a type of queue that uses a
// condition variable and a circular buffer
// BoundedQueue interface. this implementation
// is intended to be thread safe
type CondQ2 struct {
	queue []interface{}
	head     int
	tail     int
	length   int
	capacity int
	mtx sync.Mutex      // a mutex for mutual exclusion
	cvr *sync.Cond       // a condition variable for controlling mutations to the queue
}

// Put adds an element onto the tail queue
// if the queue is full, an error is returned
func (cvq *CondQ2) Put(value interface{}) error {
	// local the mutex
	cvq.cvr.L.Lock();

	// is queue full ?
	if cvq.length == cvq.capacity {
		// return an error
		e := errors.New("queue is full")
		// don't forget to unlock
		cvq.cvr.L.Unlock();
		return e;
	}

	// queue had room, add it at the tail
	cvq.queue[cvq.tail] = value
	cvq.tail = (cvq.tail+1) % cvq.capacity
	cvq.length++

	// signal a waiter if any
	cvq.cvr.Signal()

	// unlock
	cvq.cvr.L.Unlock()

	// no error
	return nil
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (cvq *CondQ2) Get() interface{} {
	// lock the mutex
	cvq.cvr.L.Lock()

	// block until a value is in the queue
	for cvq.length == 0 {
		// releast and wait
		cvq.cvr.Wait()
	}

	// at this point there is at least one item in the queue
	// remove the head
	value := cvq.queue[cvq.head]
	cvq.head = (cvq.head + 1)  % cvq.capacity
	cvq.length--

	// unlock the mutex
	cvq.cvr.L.Unlock()
	return value
}

// Try gets a value or returns an error if the queue is empty
func (cvq *CondQ2) Try() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	cvq.cvr.L.Lock()

	// is the queue empty?
	if cvq.length > 0 {
		value = cvq.queue[cvq.head]
		cvq.head = (cvq.head + 1)  % cvq.capacity
		cvq.length--
	} else {
		value = nil
		err = errors.New("queue is empty");
	}
	
	// unlock the mutex
	cvq.cvr.L.Unlock()
	return value, err
	
}

// Len is the current number of elements in the queue 
func (cvq *CondQ2) Len() int {
	return cvq.length
}

// Cap is the maximum number of elements the queue can hold
func (cvq *CondQ2) Cap() int {
	return cap(cvq.queue)
}

// String
func (cvq *CondQ2) String() string {return ""}

// NewCondQ2 is a factory for creating bounded queues
// that use a condition variable and circular buffer. It returns
// an instance of pointer to BoundedQueue
func NewCondQ2(size int) BoundedQueue {
	var cvq CondQ2
	
	// allocate the whole slice during init
	cvq.queue = make([]interface{},size,size)
	cvq.head = 0
	cvq.tail = 0
	cvq.length = 0
	cvq.capacity = size
	cvq.mtx = sync.Mutex{} // unlock mutex
	cvq.cvr = sync.NewCond(&cvq.mtx)

	return &cvq
}
