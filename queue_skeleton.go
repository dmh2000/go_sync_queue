package queue

import (
	"errors"
	"sync"
)

// Skeleton is a type of queue that uses a mutex and condition
// variable to implement the BoundedQueue interface.
// this version is a skeleton illustrating the mutual exclusion
// but has no backing data structure. It will fail the tests
type Skeleton struct {
	// -- some data structure for backing the queue
	length   int
	capacity int
	mtx sync.Mutex      // a mutex for mutual exclusion
	putcv *sync.Cond    // a condition variable for controlling Puts
	getcv *sync.Cond    // a condition variable for controlling Gets
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (skel *Skeleton) TryPut(value interface{}) error {
	// lock the mutex
	skel.putcv.L.Lock();
	defer skel.putcv.L.Unlock()

	// is queue full ?
	if skel.length == skel.capacity {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it at the tail
	// ==> enqueueing a value
	// ...

	// signal a Get to wake up
	skel.getcv.Signal()
	
	// no error
	return nil
} 

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (skel *Skeleton) Put(value interface{})  {
	// lock the mutex
	skel.putcv.L.Lock()
	defer skel.putcv.L.Unlock()


	// block until a value is in the queue
	for skel.length == skel.capacity {
		// release and wait
		skel.putcv.Wait()
	}
	
	// queue has room, add it at the tail
	// ==> enqueueing a value
	// ...


	// signal a Get to wake up
	skel.getcv.Signal()
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (skel *Skeleton) Get() interface{} {
	var value interface{}

	// lock the mutex
	skel.getcv.L.Lock()
	defer skel.getcv.L.Unlock()

	// block until a value is in the queue
	for skel.length == 0 {
		// release and wait
		skel.getcv.Wait()
	}

	// at this point there is at least one item in the queue
	// ==> dequeuing a value
	// ...
	value = 0

	// signal a Put to wake up
	skel.putcv.Signal()

	return value
}

// TryGet attempts to get a value
// if the queue is empty returns an error
func (skel *Skeleton) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	skel.getcv.L.Lock()
	defer skel.getcv.L.Unlock()

	// does the queue have elements?
	if skel.length > 0 {
		// ==> dequeuing a value
		// ...
		value = 0
	} else {
		value = nil
		err = errors.New("queue is empty");
	}

	// signal a Put to wake up
	skel.putcv.Signal()
	
	// unlock the mutex
	return value, err
}

// Len is the current number of elements in the queue 
func (skel *Skeleton) Len() int {
	return skel.length
}

// Cap is the maximum number of elements the queue can hold
func (skel *Skeleton) Cap() int {
	return skel.capacity
}

// Close handles any required cleanup
func (skel *Skeleton) Close() {
	// noop
}

// String
func (skel *Skeleton) String() string {return ""}

// NewSkeletonQueue is a factory for creating bounded queues
// that use a condition variable and circular buffer. It returns
// an instance of pointer to BoundedQueue
func NewSkeletonQueue(size int) BoundedQueue {
	var skel Skeleton
	
	// allocate the whole slice during init
	skel.length = 0
	skel.capacity = size
	skel.mtx = sync.Mutex{} 

	// both condition variables get the same mutex
	// but wakeups go from put to get and vice versa
	skel.putcv = sync.NewCond(&skel.mtx)
	skel.getcv = sync.NewCond(&skel.mtx)

	return &skel
}
