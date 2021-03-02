package queue

import "testing"

// test variables
const queueSize int = 8

// test a bounded queue
func q1(t *testing.T, q BoundedQueue) {
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
	err = q.Put(99)
	if err == nil {
		t.Error("err should be nil")
	}
	// check length is unchanged
	if q.Len() != queueSize {
		t.Error("length should == queueSize")
	}

	// remove all items
	j := queueSize - 1
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

// func TestCondvar1(t *testing.T) {
// 	// using condition variable queue
// 	q1(t,NewCondQ(queueSize))
// }

func TestCondvar2(t *testing.T) {
	// using condition variable queue
	q1(t,NewCondQ2(queueSize))
}
func TestChannel1(t *testing.T) {
	// using channel
	q1(t,NewCHQ(queueSize))
}
