package queue

import "fmt"

// BoundedQueue is a FIFO queue with a bound on the number of elements in the queue
type  BoundedQueue interface {

	// add an element onto the tail queue
	// if the queue is full, an error is returned
 	Put(value interface{}) 

	// add an element onto the tail queue
	// if the queue is full the call blocks
	TryPut(value interface{}) error

	// get an element from the head of the queue
	// if the queue is empty the get'er blocks
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
}
