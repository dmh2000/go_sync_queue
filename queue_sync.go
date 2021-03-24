package queue

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

// SynchronizedQueueImpl is an implementation of the SynchronizedQueue interface
// using a Mutex and 2 condition variables.
type SynchronizedQueueImpl struct {
	queue Queue	    // some data structure for backing the queue
	mtx sync.Mutex      // a mutex for mutual exclusion
	putcv *sync.Cond    // a condition variable for controlling Puts
	getcv *sync.Cond    // a condition variable for controlling Gets
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (sq *SynchronizedQueueImpl) TryPut(value interface{}) error {
	// lock the mutex
	sq.putcv.L.Lock();
	defer sq.putcv.L.Unlock()

	// is queue full ?
	if sq.queue.Len() == sq.queue.Cap() {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it at the tail
	// ==> enqueueing a value
	sq.queue.Push(value)

	// signal a Get to wake up
	sq.getcv.Signal()
	
	// no error
	return nil
} 

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (sq *SynchronizedQueueImpl) Put(value interface{})  {
	// lock the mutex
	sq.putcv.L.Lock()
	defer sq.putcv.L.Unlock()


	// block until a value is in the queue
	for sq.queue.Len() == sq.queue.Cap() {
		// release and wait
		sq.putcv.Wait()
	}
	
	// queue has room, add it at the tail
	// ==> enqueueing a value
	sq.queue.Push(value)


	// signal a Get to wake up
	sq.getcv.Signal()
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (sq *SynchronizedQueueImpl) Get() interface{} {
	var value interface{}

	// lock the mutex
	sq.getcv.L.Lock()
	defer sq.getcv.L.Unlock()

	// block until a value is in the queue
	for sq.queue.Len() == 0 {
		// release and wait
		sq.getcv.Wait()
	}

	// at this point there is at least one item in the queue
	// ==> dequeuing a value
	// ...
	value, err := sq.queue.Pop()
	if err != nil {
		log.Fatal(err)
	}

	// signal a Put to wake up
	sq.putcv.Signal()

	return value
}

// TryGet attempts to get a value
// if the queue is empty returns an error
func (sq *SynchronizedQueueImpl) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	sq.getcv.L.Lock()
	defer sq.getcv.L.Unlock()

	// does the queue have elements?
	if sq.queue.Len() > 0 {
		value, err = sq.queue.Pop()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		value = nil
		err = errors.New("queue is empty");
	}

	// signal a Put to wake up
	sq.putcv.Signal()
	
	// unlock the mutex
	return value, err
}

// Len is the current number of elements in the queue 
func (sq *SynchronizedQueueImpl) Len() int {
	return sq.queue.Len()
}

// Cap is the maximum number of elements the queue can hold
func (sq *SynchronizedQueueImpl) Cap() int {
	return sq.queue.Cap()
}

// Close handles any required cleanup
func (sq *SynchronizedQueueImpl) Close() {
	// noop
}

// String
func (sq *SynchronizedQueueImpl) String() string {
	var s string = "SynchronizedQueue"
	if sq.queue != nil {
		s = fmt.Sprintf("%s:%s",s,sq.queue.String())
	} else {
		s = fmt.Sprintf("%s : no queue",s);
	}
	return s
}


// NewSynchronizedQueue is a factory for creating bounded queues
// that use a mutex and condition variable
// returns an instance of SynchronizedQueue
func NewSynchronizedQueue(q Queue) SynchronizedQueue {
	var sq SynchronizedQueueImpl
	
	// attach the underlying queue data structure
	sq.queue = q

	// both condition variables get the same mutex
	// but wakeups go from put to get and vice versa
	sq.mtx = sync.Mutex{} 
	sq.putcv = sync.NewCond(&sq.mtx)
	sq.getcv = sync.NewCond(&sq.mtx)

	return &sq
}
