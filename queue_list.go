package queue

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

// ListQ is a type of queue that uses a
// condition variable and lists to implement the
// BoundedQueue interface. this implementation
// is intended to be thread safe
type ListQ struct {
	queue *list.List	 // contains the elements currently in the queue
	capacity int	     // maximum number of elements the queue can hold
	mtx sync.Mutex       // a mutex for mutual exclusion
	getcv *sync.Cond     // a condition variable for controlling gets
	putcv *sync.Cond       // a condition variable for controlling puts
}


// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (lsq *ListQ) Put(value interface{})  {
	// lock the mutex
	lsq.putcv.L.Lock()
	defer lsq.putcv.L.Unlock()


	// block until a value is in the queue
	for lsq.queue.Len() == lsq.capacity {
		// release and wait
		lsq.putcv.Wait()
	}
	
	// queue had room, add it
	lsq.queue.PushBack(value)

	// signal a Get to wakeup
	lsq.getcv.Signal()
} 

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (lsq *ListQ) TryPut(value interface{}) error {
	// lock the mutex
	lsq.putcv.L.Lock();
	defer lsq.putcv.L.Unlock()

	// is queue full ?
	if lsq.queue.Len() == lsq.capacity {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it
	lsq.queue.PushBack(value)

	// signal a Get if any
	lsq.getcv.Signal()

	// no error
	return nil
} 

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (lsq *ListQ) Get() interface{} {
	// lock the mutex
	lsq.getcv.L.Lock()
	defer lsq.getcv.L.Unlock()

	// block until a value is in the queue
	for lsq.queue.Len() == 0 {
		// release and wait
		lsq.getcv.Wait()
	}

	// at this point there is at least one item in the queue
	// remove the head
	value := lsq.queue.Remove(lsq.queue.Front())

	// signal a Put to wakeup
	lsq.putcv.Signal()

	return value
}

// TryGet gets a value or returns an error if the queue is empty
func (lsq *ListQ) TryGet() (interface{}, error) {
	var value interface{}
	var err error

	// lock the mutex
	lsq.getcv.L.Lock()
	defer lsq.getcv.L.Unlock()

	// is the queue empty?
	if lsq.queue.Len() > 0 {
		value = lsq.queue.Remove(lsq.queue.Front())
	} else {
		value = nil
		err = errors.New("queue is empty");
	}
	
	// signal a Put to wakeup
	lsq.putcv.Signal()

	return value, err
}

// Len is the current number of elements in the queue 
func (lsq *ListQ) Len() int {
	return lsq.queue.Len()
}

// Cap is the maximum number of elements the queue can hold
func (lsq *ListQ) Cap() int {
	return lsq.capacity;
}

// Close handles any required cleanup
func (lsq *ListQ) Close() {
	// noop
}

// String
func (lsq *ListQ) String() string {
	return fmt.Sprintf("List Len:%v Cap:%v",lsq.Len(),lsq.Cap())
}

// NewListQueue is a factory for creating bounded queues
// that use a condition variable and lists. It returns
// an instance of pointer to BoundedQueue
func NewListQueue(size int) BoundedQueue {
	var lsq ListQ

	lsq.capacity = size
	lsq.queue = list.New()
	lsq.mtx = sync.Mutex{} 
	lsq.putcv = sync.NewCond(&lsq.mtx)
	lsq.getcv = sync.NewCond(&lsq.mtx)

	return &lsq
}

