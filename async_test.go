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
func producer1(q SynchronizedQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}

	// cleanup 
	// for channels, this closes it for further Puts
	// any remaing data is still available for Gets
	// it is a noop for the mutex/condition variable methods
	q.Close()

	// mark it done
	wg.Done()
}

func consumer1(q SynchronizedQueue, t *testing.T, wg *sync.WaitGroup)  {
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
func producer2(q *NativeIntQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}
	wg.Done()
}

func consumer2(q *NativeIntQueue, t *testing.T, wg *sync.WaitGroup) {
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
func producer3(q SynchronizedQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
		q.Put(i)
	}

	// cleanup 
	// for channels, this closes it for further Puts
	// any remaing data is still available for Gets
	// it is a noop for the mutex/condition variable methods
	q.Close()
	
	// mark it done
	wg.Done()
}

func consumer3(q SynchronizedQueue, t *testing.T, wg *sync.WaitGroup)  {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
		value := q.Get()
		// convert to int
		v := value.(int)
		if v != i {
			t.Error("v should == i",v,i)
		}
	}
	wg.Done()
}


// 3 - blocking with random time delays
func producer4(q *NativeIntQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
		q.Put(i)
	}

	// cleanup 
	// for channels, this closes it for further Puts
	// any remaing data is still available for Gets
	// it is a noop for the mutex/condition variable methods
	q.Close()
	
	// mark it done
	wg.Done()
}

func consumer4(q *NativeIntQueue, t *testing.T, wg *sync.WaitGroup)  {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
		value := q.Get()
		// convert to int
		if value != i {
			t.Error("v should == i",value,i)
		}
	}
	wg.Done()
}

func async1(t *testing.T, q SynchronizedQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer1(q,&wg)
	go consumer1(q,t,&wg)
	wg.Wait()
}

func async2(t *testing.T, q *NativeIntQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer2(q, &wg)
	go consumer2(q,t,&wg)
	wg.Wait()
}

func async3(t *testing.T, q SynchronizedQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer3(q, &wg)
	go consumer3(q,t,&wg)
	wg.Wait()
}

func async4(t *testing.T, q *NativeIntQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer4(q, &wg)
	go consumer4(q,t,&wg)
	wg.Wait()
}

func TestChannelAsync(t *testing.T) {
	async1(t,NewChannelQueue(aqsize))
	async3(t,NewChannelQueue(aqsize))
}

func TestListAsync(t *testing.T) {
	async1(t,NewSyncList(aqsize))
	async3(t,NewSyncList(aqsize))
}
func TestCircularAsync(t *testing.T) {
	async1(t,NewSyncCircular(aqsize))
	async3(t,NewSyncCircular(aqsize))
}

func TestRingAsync(t *testing.T) {
	async1(t,NewSyncRing(aqsize))
	async3(t,NewSyncRing(aqsize))
}

func TestComboAsync(t *testing.T) {
	async1(t,NewSyncCircular(aqsize))
	async3(t,NewSyncList(aqsize))
}

func TestNativeAsync(t *testing.T) {
	async2(t,NewNativeQueue(aqsize))
	async4(t,NewNativeQueue(aqsize))
}