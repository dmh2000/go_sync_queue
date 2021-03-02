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


var bqsize int = 100000

// func BenchmarkCondvar1(b *testing.B) {
// 	// using condition variable queue (list)
// 	b1(b,NewCondQ(bqsize))
// }

func BenchmarkCondvar2(b *testing.B) {
	// using condition variable queue (circular buffer)
	b1(b,NewCondQ2(bqsize))
}

func BenchmarkChannel1(b *testing.B) {
	// using channel
	b1(b,NewCHQ(bqsize))
}
