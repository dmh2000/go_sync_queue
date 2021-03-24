package queue

import "fmt"

// Queue - interface for a simple, non-thread-safe queue
type Queue interface {
	// current number of elements in the queue
	Len() int

	// maximum number of elements allowed in queue
	Cap() int

	// enqueue a value on the tail of the queue
	Push(value interface{})  error

	// dequeue and return a value from the head of the queue
	Pop() (interface{},error)

	// string representation
	fmt.Stringer
}

// SynchronizedQueue is a queue with a bound on the number of elements in the queue
// this interface does not promise thread-safety
type  SynchronizedQueue interface {

	// add an element onto the tail queue
	// if the queue is full, the caller blocks
 	Put(value interface{}) 

	// add an element onto the tail queue
	// if the queue is full an error is returned
	TryPut(value interface{}) error

	// get an element from the head of the queue
	// if the queue is empty the caller blocks
	Get() interface{}

	// try to get an element from the head of the queue
	// if the queue is empty an error is returned
	TryGet() (interface{}, error)

	// current number of elements in the queue 
 	Len() int

	// capacity maximum number of elements the queue can hold
	Cap() int

	// close any resources (required for channel version)
	Close()
	
	// string representation
	fmt.Stringer
}

