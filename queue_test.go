package queue

import (
	"fmt"
	"sync"
	"testing"
)

// test variables
const queueSize int = 8

// test an instance of a BoundedQueue
func sync1(t *testing.T, q BoundedQueue) {
	var err error

	if q == nil {
		t.Error("q should not be nil")
	}

	// check length
	if q.Len() != 0 {
		t.Error("length should be 0",q.Len())
	}

	// check capacity
	if q.Cap() != queueSize {
		t.Error("capacity should == queueSize",q.Cap(),queueSize)
	}

	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.TryPut(i)
		//length should be == i at this point
		if q.Len() != (i+1) {
			t.Error("length should == i+1",q.Len(),i+1)
		}
	}

	// check the length, should be == capacity
	if q.Len() != q.Cap() {
		t.Error("length should == capacity")
	}

	// try to add one more
	err = q.TryPut(99)
	if err == nil {
		t.Error("err should be nil")
	}
	// check length is unchanged
	if q.Len() != queueSize {
		t.Error("length should == queueSize")
	}

	// remove all items
	j := q.Len() - 1
	for i:=0;i<q.Cap();i++ {
		v := q.Get().(int)
		if v != i {
			t.Error("v should == i",v,i)
		}
		// length should decrease
		if q.Len() != j {
			t.Error("length should == i",q.Len(),j)
		}
		j--
	}

	// check the length == 0
	if q.Len() != 0 {
		t.Error("length should == 0",q.Len())
	}
}


// test an instance of a BoundedQueue
func sync2(t *testing.T, q *NativeIntQ) {
	var err error

	if q == nil {
		t.Error("q should not be nil")
	}

	// check length
	if q.Len() != 0 {
		t.Error("length should be 0",q.Len())
	}

	// check capacity
	if q.Cap() != queueSize {
		t.Error("capacity should == queueSize",q.Cap(),queueSize)
	}

	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
		//length should be == i at this point
		if q.Len() != (i+1) {
			t.Error("length should == i+1",q.Len(),i+1)
		}
	}

	// check the length, should be == capacity
	if q.Len() != q.Cap() {
		t.Error("length should == capacity")
	}

	// try to add one more
	err = q.TryPut(99)
	if err == nil {
		t.Error("err should be nil")
	}
	// check length is unchanged
	if q.Len() != queueSize {
		t.Error("length should == queueSize")
	}

	// remove all items
	j := q.Len() - 1
	for i:=0;i<q.Cap();i++ {
		v := q.Get()
		if v != i {
			t.Error("v should == i",v,i)
		}
		// length should decrease
		if q.Len() != j {
			t.Error("length should == i",q.Len(),j)
		}
		j--
	}

	// check the length == 0
	if q.Len() != 0 {
		t.Error("length should == 0",q.Len())
	}
}

// ====================
// SYNCHRONOUS TESTS
// ====================

// CHANNEL
func TestChannelSync(t *testing.T) {
	// using channel
	sync1(t,NewChannelQueue(queueSize))
}

// LIST
func TestListSync(t *testing.T) {
	// using condition variable queue
	sync1(t,NewListQueue(queueSize))
}

// CIRCULAR BUFFER
func TestCircularSync(t *testing.T) {
	// using condition variable queue
	sync1(t,NewCircularQueue(queueSize))
}

// NATIVE QUEUE
func TestQueueNativeSync(t *testing.T) {
	// using condition variable queue
	sync2(t,NewNativeQueue(queueSize))
}

// Strings
func TestQueueStrings(t *testing.T) {
	fmt.Println(NewChannelQueue(queueSize))
	fmt.Println(NewListQueue(queueSize))
	fmt.Println(NewCircularQueue(queueSize))
	fmt.Println(NewNativeQueue(queueSize))
}

// ==================
// Asynchronous Tests
// ==================
func producer1(q BoundedQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}

	wg.Done()
}

func consumer1(q BoundedQueue, t *testing.T, wg *sync.WaitGroup)  {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		value := q.Get()
		// convert to int
		v := value.(int)
		if v != i {
			t.Error("v should == i",v,i)
		}
	}
	wg.Done()
}

func producer2(q *NativeIntQ, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}
	wg.Done()
}

func consumer2(q *NativeIntQ, t *testing.T, wg *sync.WaitGroup) {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		v := q.Get()
		if v != i {
			t.Error("v should == i",v,i)
		}
	}
	wg.Done()
}

func async1(t *testing.T, q BoundedQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer1(q,&wg)
	go consumer1(q,t,&wg)
	wg.Wait()
}

func async2(t *testing.T, q *NativeIntQ) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer2(q, &wg)
	go consumer2(q,t,&wg)
	wg.Wait()
}

func TestChannelAsync(t *testing.T) {
	async1(t,NewChannelQueue(queueSize))
}

func TestListAsync(t *testing.T) {
	async1(t,NewListQueue(queueSize))
}
func TestCircularAsync(t *testing.T) {
	async1(t,NewCircularQueue(queueSize))
}

func TestNativeAsync(t *testing.T) {
	async2(t,NewNativeQueue(queueSize))
}