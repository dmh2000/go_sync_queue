package queue

import (
	"sync"
	"testing"
)

//  go test -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out

// queue size
var bqsize int = 20


func b1(b *testing.B, q SynchronizedQueue) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		var x interface{}
		x = i

		q.Put(x)
		//length should be == i at this point
		if q.Len() != (i+1) {
			b.Error("length should == i+1",q.Len(),i+1)
		}
	}

	// remove all items
	for i:=0;i<q.Cap();i++ {
		v := q.Get().(int)
		if v != i {
			b.Error("v should == i",v,i)
		}
	}
}

func b2(b *testing.B, q *NativeIntQueue) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		var x int
		x = i

		q.Put(x)
		//length should be == i at this point
		if q.Len() != (i+1) {
			b.Error("length should == i+1",q.Len(),i+1)
		}
	}

	// remove all items
	for i:=0;i<q.Cap();i++ {
		v := q.Get()
		if v != i {
			b.Error("v should == i",v,i)
		}
	}
}

// Synchronous Benchmarks

func BenchmarkQueueChannelSync(b *testing.B) {
	// using channel
	for i:=0;i<b.N;i++ {
		b1(b,NewChannelQueue(bqsize))
	}
}

func BenchmarkListSync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		b1(b,NewSyncList(bqsize))
	}
}

func BenchmarkCircularSync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		b1(b,NewSyncCircular(bqsize))
	}
}

func BenchmarkRingSync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		b1(b,NewSyncRing(bqsize))
	}
}

func BenchmarkSliceSync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		b1(b,NewSyncSlice(bqsize))
	}
}

func BenchmarkQueueNativeSync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		b2(b,NewNativeQueue(bqsize))
	}
}

// Asynchronous benchmarks


// ==================
// Asynchronous Tests
// ==================

// - blocking with no delays
func producer1a(q SynchronizedQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}

	wg.Done()
}

func consumer1a(q SynchronizedQueue, b *testing.B, wg *sync.WaitGroup)  {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		value := q.Get()
		// convert to int
		v := value.(int)
		if v != i {
			b.Error("v should == i",v,i)
		}
	}
	wg.Done()
}

// - blocking with native ints, no delays
func producer2a(q *NativeIntQueue, wg *sync.WaitGroup) {
	// fill the queue with ints
	for i:=0;i<q.Cap();i++ {
		q.Put(i)
	}
	wg.Done()
}

func consumer2a(q *NativeIntQueue, b *testing.B, wg *sync.WaitGroup) {
	// consume all items
	for i:=0;i<q.Cap();i++ {
		v := q.Get()
		if v != i {
			b.Error("v should == i",v,i)
		}
	}
	wg.Done()
}


func asyncb1(b *testing.B, q SynchronizedQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer1a(q,&wg)
	go consumer1a(q,b,&wg)
	wg.Wait()
}

func asyncb2(b *testing.B, q *NativeIntQueue) {
	var wg sync.WaitGroup

	wg.Add(2)
	go producer2a(q,&wg)
	go consumer2a(q,b,&wg)
	wg.Wait()
}

func BenchmarkQueueChannelAsync(b *testing.B) {
	// using channel
	for i:=0;i<b.N;i++ {
		asyncb1(b,NewChannelQueue(bqsize))
	}
}

func BenchmarkListAsync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		asyncb1(b,NewSyncList(bqsize))
	}
}

func BenchmarkCircularAsync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		asyncb1(b,NewSyncCircular(bqsize))
	}
}

func BenchmarkRingAsync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		asyncb1(b,NewSyncRing(bqsize))
	}
}

func BenchmarkSliceAsync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		asyncb1(b,NewSyncSlice(bqsize))
	}
}


func BenchmarkQueueNativeAsync(b *testing.B) {
	// using condition variable queue
	for i:=0;i<b.N;i++ {
		asyncb2(b,NewNativeQueue(bqsize))
	}
}

