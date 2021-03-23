package queue

import (
	"errors"
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
func (bq *SynchronizedQueueImpl) TryPut(value interface{}) error {
	// lock the mutex
	bq.putcv.L.Lock();
	defer bq.putcv.L.Unlock()

	// is queue full ?
	if bq.queue.Len() == bq.queue.Cap() {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it at the tail
	// ==> enqueueing a value
	bq.queue.Push(value)

	// signal a Get to wake up
	bq.getcv.Signal()
	
	// no error
	return nil
} 

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (bq *SynchronizedQueueImpl) Put(value interface{})  {
	// lock the mutex
	bq.putcv.L.Lock()
	defer bq.putcv.L.Unlock()


	// block until a value is in the queue
	for bq.queue.Len() == bq.queue.Cap() {
		// release and wait
		bq.putcv.Wait()
	}
	
	// queue has room, add it at the tail
	// ==> enqueueing a value
	bq.queue.Push(value)


	// signal a Get to wake up
	bq.getcv.Signal()
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (bq *SynchronizedQueueImpl) Get() interface{} {
	var value interface{}

	// lock the mutex
	bq.getcv.L.Lock()
	defer bq.getcv.L.Unlock()

	// block until a value is in the queue
	for bq.queue.Len() == 0 {
		// release and wait
		bq.getcv.Wait()
	}

	// at this point there is at least one item in the queue
	// ==> dequeuing a value
	// ...
	value, err := bq.queue.Pop()
	if err != nil {
		log.Fatal(err)
	}

	// signal a Put to wake up
	bq.putcv.Signal()

	return value
}

// TryGet attempts to get a value
// if the queue is empty returns an error
func (bq *SynchronizedQueueImpl) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	bq.getcv.L.Lock()
	defer bq.getcv.L.Unlock()

	// does the queue have elements?
	if bq.queue.Len() > 0 {
		value, err = bq.queue.Pop()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		value = nil
		err = errors.New("queue is empty");
	}

	// signal a Put to wake up
	bq.putcv.Signal()
	
	// unlock the mutex
	return value, err
}

// Len is the current number of elements in the queue 
func (bq *SynchronizedQueueImpl) Len() int {
	return bq.queue.Len()
}

// Cap is the maximum number of elements the queue can hold
func (bq *SynchronizedQueueImpl) Cap() int {
	return bq.queue.Cap()
}

// Close handles any required cleanup
func (bq *SynchronizedQueueImpl) Close() {
	// noop
}

// String
func (bq *SynchronizedQueueImpl) String() string {return ""}


// NewSynchronizedQueue is a factory for creating bounded queues
// that use a mutex and condition variable
// returns an instance of SynchronizedQueue
func NewSynchronizedQueue(q Queue) SynchronizedQueue {
	var bq SynchronizedQueueImpl
	
	// attach the underlying queue data structure
	bq.queue = q

	// both condition variables get the same mutex
	// but wakeups go from put to get and vice versa
	bq.mtx = sync.Mutex{} 
	bq.putcv = sync.NewCond(&bq.mtx)
	bq.getcv = sync.NewCond(&bq.mtx)

	return &bq
}
