package queue

import "testing"

func b1(b *testing.B, q BoundedQueue) {

	for k:=0;k<b.N;k++ {
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
}

func b2(b *testing.B, q *NativeInt) {

	for k:=0;k<b.N;k++ {
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
}

var bqsize int = 1000000


func BenchmarkQueueChannel(b *testing.B) {
	// using channel
	b1(b,NewChannelQueue(queueSize))
}

func BenchmarkList(b *testing.B) {
	// using condition variable queue
	b1(b,NewListQueue(queueSize))
}

func BenchmarkCircular(b *testing.B) {
	// using condition variable queue
	b1(b,NewCircularQueue(queueSize))
}

func BenchmarkQueueNative(b *testing.B) {
	// using condition variable queue
	b2(b,NewNativeQueue(queueSize))
}

