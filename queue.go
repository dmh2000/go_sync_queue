package queue

import "fmt"

// BoundedQueue is a FIFO queue with a bound on the number of elements in the queue
type  BoundedQueue interface {

	// add an element onto the tail queue
	// if the queue is full, an error is returned
	// this is in contrast to an implementation that would have the put'er
	// block until room in the queue is available
 	Put(value interface{}) error 

	// get an element from the head of the queue
	// if the queue is empty the get'er blocks
	Get() interface{}

	// try to get an element from the head of the queue
	// if the queue is empty an error is returned
	Try() (interface{}, error)

	// current number of elements in the queue 
 	Len() int

	// capacity maximum number of elements the queue can hold
	Cap() int

	// string representation
	fmt.Stringer
}