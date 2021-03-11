package queue

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// test variables
const aqsize int = 8


// ==================
// Asynchronous Tests
// ==================

// - blocking with no delays
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

// - blocking with native ints, no delays
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

// 3 - blocking with random time delays
func producer3(q BoundedQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(250)) * time.Millisecond)
		q.Put(i)
	}

	wg.Done()
}

func consumer3(q BoundedQueue, t *testing.T, wg *sync.WaitGroup)  {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(250)) * time.Millisecond)
		value := q.Get()
		// convert to int
		v := value.(int)
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

func async3(t *testing.T, q BoundedQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer3(q, &wg)
	go consumer3(q,t,&wg)
	wg.Wait()
}


func TestChannelAsync(t *testing.T) {
	async1(t,NewChannelQueue(aqsize))
	async3(t,NewChannelQueue(aqsize))
}

func TestListAsync(t *testing.T) {
	async1(t,NewListQueue(aqsize))
	async3(t,NewListQueue(aqsize))
}
func TestCircularAsync(t *testing.T) {
	async1(t,NewCircularQueue(aqsize))
	async3(t,NewCircularQueue(aqsize))
}

func TestNativeAsync(t *testing.T) {
	async2(t,NewNativeQueue(aqsize))
}