package queue

import (
	"errors"
	"fmt"
)

// ChannelQ is a type of queue that uses a
// condition variable and lists to implement the
// SynchronizedQueue interface. this implementation
// is intended to be thread safe
type ChannelQ struct {
	channel chan interface{} // buffered channel with specified capacity 
}

// TryPut adds an element onto the tail queue
// if the queue is full, an error is returned
func (chq *ChannelQ) TryPut(value interface{}) error {
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
func (chq *ChannelQ) Put(value interface{}) {
	chq.channel <- value
}

// Get returns an element from the head of the queue
// if the queue is empty,the caller blocks
func (chq *ChannelQ) Get() interface{} {
	// get a value or block
	return <-chq.channel
}

// TryGet gets a value or returns an error if the queue is empty
func (chq *ChannelQ) TryGet() (interface{}, error) {
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
func (chq *ChannelQ) Len() int {
	return len(chq.channel)
}

// Cap is the maximum number of elements the queue can hold
func (chq *ChannelQ) Cap() int {
	return cap(chq.channel)
}

// Close required to close the channel so it doesn't leak
// should only be called by Put threads
func (chq *ChannelQ) Close() {
	close(chq.channel)
}

// String
func (chq *ChannelQ) String() string {
	return fmt.Sprintf("ChannelQ Len:%v Cap:%v",chq.Len(),chq.Cap())
}

// NewChannelQueue is a factory for creating bounded queues
// that uses a channel
// It returns an instance of pointer to SynchronizedQueue
func NewChannelQueue(size int) SynchronizedQueue {
	var chq ChannelQ

	chq.channel = make(chan interface{}, size)

	return &chq
}
