package queue

import (
	"errors"
)

// Channel is a type of queue that uses a
// condition variable and lists to implement the
// BoundedQueue interface. this implementation
// is intended to be thread safe
type Channel struct {
	channel chan interface{} // buffered channel with specified capacity 
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (chq *Channel) TryPut(value interface{}) error {
	var err error

	err = nil

	// attempt to insert the value into the buffered channel
	select {
	// send it if there is room
	case chq.channel <- value:
		// no action
	default:
		// couldn't send, buffered channel is full
		err = errors.New("queue is full")
	}

	return err
} 

// Put adds an element to the tail of the queue
// if the queue is full the function blocks
func (chq *Channel) Put(value interface{}) {
	chq.channel <- value
}

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (chq *Channel) Get() interface{} {
	// get a value or block
	return <-chq.channel
}

// TryGet gets a value or returns an error if the queue is empty
func (chq *Channel) TryGet() (interface{}, error) {
	var err error
	var value interface{}

	value = nil
	select {
	case value = <-chq.channel:
		// no action
	default:
		err = errors.New("queue is empty")
	}
	
	return value,err
}

// Len is the current number of elements in the queue 
func (chq *Channel) Len() int {
	return len(chq.channel)
}

// Cap is the maximum number of elements the queue can hold
func (chq *Channel) Cap() int {
	return cap(chq.channel)
}

// String
func (chq *Channel) String() string {return ""}

// NewChannelQueue is a factory for creating bounded queues
// that uses a channel
// It returns an instance of pointer to BoundedQueue
func NewChannelQueue(size int) BoundedQueue {
	var chq Channel

	chq.channel = make(chan interface{}, size)

	return &chq
}
