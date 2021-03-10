package queue

import (
	"errors"
	"sync"
)

// NativeInt is a type of queue that uses a
// condition variable and a circular buffer
// BoundedQueue interface. this implementation
// is intended to be thread safe
type NativeInt struct {
	queue []int
	head     int
	tail     int
	length   int
	capacity int
	mtx sync.Mutex      // a mutex for mutual exclusion
	cvr *sync.Cond       // a condition variable for controlling mutations to the queue
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (nvq *NativeInt) TryPut(value int) error {
	// local the mutex
	nvq.cvr.L.Lock();
	defer nvq.cvr.L.Unlock()

	// is queue full ?
	if nvq.length == nvq.capacity {
		// return an error
		e := errors.New("queue is full")
		return e;
	}

	// queue had room, add it at the tail
	nvq.queue[nvq.tail] = value
	nvq.tail = (nvq.tail+1) % nvq.capacity
	nvq.length++

	// signal a waiter if any
	nvq.cvr.Signal()

	// no error
	return nil
} 

// Put adds an element onto the tail queue
// if the queue is full the function blocks
func (nvq *NativeInt) Put(value int)  {
	// local the mutex
	nvq.cvr.L.Lock()
	defer nvq.cvr.L.Unlock()


	// block until a value is in the queue
	for nvq.length == nvq.capacity {
		// releast and wait
		nvq.cvr.Wait()
	}
	
	// queue has room, add it at the tail
	nvq.queue[nvq.tail] = value
	nvq.tail = (nvq.tail+1) % nvq.capacity
	nvq.length++

	// signal a waiter if any
	nvq.cvr.Signal()
} 


// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (nvq *NativeInt) Get() int {
	// lock the mutex
	nvq.cvr.L.Lock()
	defer nvq.cvr.L.Unlock()

	// block until a value is in the queue
	for nvq.length == 0 {
		// releast and wait
		nvq.cvr.Wait()
	}

	// at this point there is at least one item in the queue
	// remove the head
	value := nvq.queue[nvq.head]
	nvq.head = (nvq.head + 1)  % nvq.capacity
	nvq.length--

	return value
}

// Try gets a value or returns an error if the queue is empty
func (nvq *NativeInt	) Try() (int, error) {
	var value int
	var err error

	// lock the mutex
	nvq.cvr.L.Lock()
	defer nvq.cvr.L.Unlock()

	// is the queue empty?
	if nvq.length > 0 {
		value = nvq.queue[nvq.head]
		nvq.head = (nvq.head + 1)  % nvq.capacity
		nvq.length--
	} else {
		value = 0
		err = errors.New("queue is empty");
	}
	
	return value, err
	
}

// Len is the current number of elements in the queue 
func (nvq *NativeInt	) Len() int {
	return nvq.length
}

// Cap is the maximum number of elements the queue can hold
func (nvq *NativeInt	) Cap() int {
	return cap(nvq.queue)
}

// String
func (nvq *NativeInt	) String() string {return ""}

// NewNativeQueue is a factory for creating bounded queues
// that use a condition variable and circular buffer. It returns
// an instance of pointer to BoundedQueue
func NewNativeQueue(size int) *NativeInt {
	var nvq NativeInt
	
	// allocate the whole slice during init
	nvq.queue = make([]int,size,size)
	nvq.head = 0
	nvq.tail = 0
	nvq.length = 0
	nvq.capacity = size
	nvq.mtx = sync.Mutex{} // unlock mutex
	nvq.cvr = sync.NewCond(&nvq.mtx)

	return &nvq
}
