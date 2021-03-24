package queue

import (
	"math/rand"
	"sort"
	"testing"
)

// test variables
const hqsize int = 8


// test an instance of a SynchronizedQueue
func heap1(t *testing.T, q SynchronizedQueue) {
	if q == nil {
		t.Error("q should not be nil")
	}

	// check length
	if q.Len() != 0 {
		t.Error("length should be 0",q.Len())
	}

	// check capacity
	if q.Cap() != hqsize {
		t.Error("capacity should == hqsize",q.Cap(),hqsize)
	}

	// create a list of random positive values
	rval := make(PrioritySlice,q.Cap())
	for i:=0;i<len(rval);i++ {
		j := rand.Intn(100)
		rval[i].value = j*10
		rval[i].priority = j
	}
	
	// insert it the priority queue
	for i:=0;i<q.Cap();i++ {
		q.TryPut(rval[i])
		//length should be == i at this point
		if q.Len() != (i+1) {
			t.Error("length should == i+1",q.Len(),i+1)
		}
		// print it 
		// fmt.Println(rval[i])
	}

	// check the length, should be == capacity
	if q.Len() != q.Cap() {
		t.Error("length should == capacity")
	}

	// cleanup 
	// for channels, this closes it for further Puts
	// any remaing data is still available for Gets
	// it is a noop for the mutex/condition variable methods
	q.Close()

	// fmt.Println()

	// sort the rval list
	sort.Sort(rval)

	// remove all items. no need to block
	// should come out in ascending order
	j := q.Len() - 1
	for i:=0;i<q.Cap();i++ {
		value, err := q.TryGet()
		if err != nil {
			t.Error(err)
		}
		// assert convert to int
		v := value.(PriorityItem)
		if v.priority != rval[i].priority {
			t.Error("heap not in order")
		}
		// fmt.Println(v)
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

// PRIORITY QUEUE using SynchronizedQueue wrapper
func TestPrioritySync(t *testing.T) {
	// test for heap to check priority  
	heap1(t,NewSyncPriority(hqsize))
}
