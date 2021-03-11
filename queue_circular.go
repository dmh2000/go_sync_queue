package queue

import (
	"errors"
	"fmt"
	"sync"
)

// CircularQ is a type of queue that uses a
// condition variable and a circular buffer
// BoundedQueue interface. this implementation
// is intended to be thread safe
type CircularQ struct {
	queue []interface{}
	head     int
	tail     int
	length   int
	capacity int
	mtx sync.Mutex       // a mutex for mutual exclusion
	putcv *sync.Cond   // a condition variable for controlling mutations to the queue
	getcv *sync.Cond
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (cir *CircularQ) TryPut(value interface{}) error {
	// lock the mutex
	cir.putcv.L.Lock();
	defer cir.putcv.L.Unlock()

	// is queue full ?
	if cir.length == cir.capacity {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it at the tail
	cir.queue[cir.tail] = value
	cir.tail = (cir.tail+1) % cir.capacity
	cir.length++

	// signal a Get to wake up
	cir.getcv.Signal()
	
	// no error
	return nil
} 

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (cir *CircularQ) Put(value interface{})  {
	// lock the mutex
	cir.putcv.L.Lock()
	defer cir.putcv.L.Unlock()


	// block until a value is in the queue
	for cir.length == cir.capacity {
		// releast and wait
		cir.putcv.Wait()
	}
	
	// queue has room, add it at the tail
	cir.queue[cir.tail] = value
	cir.tail = (cir.tail+1) % cir.capacity
	cir.length++

	// signal a get to wakeup
	cir.getcv.Signal()
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (cir *CircularQ) Get() interface{} {
	// lock the mutex
	cir.getcv.L.Lock()
	defer cir.getcv.L.Unlock()

	// block until a value is in the queue
	for cir.length == 0 {
		// releast and wait
		cir.getcv.Wait()
	}

	// at this point there is at least one item in the queue
	// remove the head
	value := cir.queue[cir.head]
	cir.head = (cir.head + 1)  % cir.capacity
	cir.length--

	// signal a put to wakeup
	cir.putcv.Signal()

	return value
}

// TryGet attempts to get a value
// if the queue is empty returns an error
func (cir *CircularQ) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	cir.getcv.L.Lock()
	defer cir.getcv.L.Unlock()

	// is the queue empty?
	if cir.length > 0 {
		value = cir.queue[cir.head]
		cir.head = (cir.head + 1)  % cir.capacity
		cir.length--
	} else {
		value = nil
		err = errors.New("queue is empty");
	}

	// signal a put to wakeup
	cir.putcv.Signal()
	
	// unlock the mutex
	return value, err
}

// Len is the current number of elements in the queue 
func (cir *CircularQ) Len() int {
	return cir.length
}

// Cap is the maximum number of elements the queue can hold
func (cir *CircularQ) Cap() int {
	return cap(cir.queue)
}

// String
func (cir *CircularQ) String() string {
	return fmt.Sprintf("Circular Len:%v Cap:%v",cir.Len(),cir.Cap())
}

// NewCircularQueue is a factory for creating bounded queues
// that use a condition variable and circular buffer. It returns
// an instance of pointer to BoundedQueue
func NewCircularQueue(size int) BoundedQueue {
	var cir CircularQ
	
	// allocate the whole slice during init
	cir.queue = make([]interface{},size,size)
	cir.head = 0
	cir.tail = 0
	cir.length = 0
	cir.capacity = size
	cir.mtx = sync.Mutex{} 
	cir.getcv = sync.NewCond(&cir.mtx)
	cir.putcv = sync.NewCond(&cir.mtx)

	return &cir
}
