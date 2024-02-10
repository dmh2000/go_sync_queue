---
title: "Synchronized Queues in Golang"
date: 2021-03-12
slug: "/syncqueue"
---

![test](https://github.com/dmh2000/go_sync_queue/actions/workflows/go.yml/badge.svg)

## Synchronized Queue in Golang

All code files are in Github at [dmh2000/go_sync_queue](https://github.com/dmh2000/go_sync_queue)

[There are eight million stories in the city. This has been one of them.](<https://en.wikipedia.org/wiki/Naked_City_(TV_series)>).

**There are eight million examples of queues in Go. This is one of them.**

### Queue

As most everyone knows, a queue is a data structure with first-in/first-out (FIFO) semantics.
A generic queue typically has more or less the following methods:

- init : initialize or otherwise create a queue
- put : adds an element to the tail of the queue.
- get : takes the element from the head of the queue.
- length : how many elements are in the queue
- capacity : how many elements the queue can hold

A bounded queue is one that also has a 'capacity'. In that case, the queue is set up to contain a fixed number of elements.

Queues are used in various algorithms, such as breadth first search. Another common use of a queue is to provide a means of communication between threads. If a queue is used for that purpose, then it may have blocking semantics. The 'get' method may block if the queue is empty. The 'put' method may block if the queue is full. If that is how it works, then there can be a TryPut method that adds an element if there is room in the queue or returns an indication that the queue is full, without blocking. Likewise there can be a TryGet method that returns an element or an indication the queue is empty. Of course in Go you have buffered channels, which are effectively queues with blocking and non-blocking support.

Typically, if a queue is used for communication between threads it is expected to live for the life of the program. It gets tricky if interthread queues are to be created and discarded over the life of execution. This can lead to memory leaks and threads hung on blocking Get's or Put's if care is not taken to make sure all goroutines are released from the queue before sending it go garbage collection. The code here doesn't detect or handle that case. It is up to the application to clean up.

#### Why have a bound?

Its possible to have an unbounded queue that never gets full. Using dynamic allocation you could implement an unbounded queue. In C++ a list or vector can grow until there is no available memory. In Go, the container/list object has no bound on its number of elements. These are unbounded until you run out of memory. There is always an implicit bound. You might want to bound a queue as a type of throttling. Say a system has processing power that can only support N things going on a time. An input queue of capacity N would give it a way to stop accepting things until less than N are working. Or, maybe the design wants to know it has enough memory to support all its functions. In many real time systems there is a practice that calls for allocating all resources during initialization so that all the constraints are known and can be analyzed.

### Queue API

Here is an interface definition for a simple FIFO queue with a set bound. This interface does not support blocking semantics. Its intended to specify the methods needed to encapsulate an underlying data structure.

```go
// Queue - interface for a simple, non-thread-safe queue
type Queue interface {
	// current number of elements in the queue
	Len() int

	// maximum number of elements allowed in queue
	Cap() int

	// enqueue a value on the tail of the queue
	Push(value interface{})  error

	// dequeue and return a value from the head of the queue
	Pop() (interface{},error)

	// string representation
	fmt.Stringer
}
```

#### Implementations

I have several implementations of the Queue interface.

- SliceQueue
  - queue using a slice of interface{}
  - data is not preallocated
  - [queue_slice.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_slice.go).
- ListQueue
  - queue using a container/list
  - [queue_list.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_list.go).
- RingQueue
  - queue using a container/ring
  - [queue_ring.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_ring.go).
- CircularQueue
  - queue using a homegrown circular buffer with preallocation
  - [queue_circular.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_circular.go).
- PriorityQueue
  - queue using a container/heap
  - [queue_priority.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_priority.go).
  - in this case, its implemented as a priority queue rather than FIFO
  - data elements have to be PriorityItem

The data elements are interface{} so any type can be used. This matches some of the approaches in the standard library for certain data structures. These implementations can be passed to any function needed a Queue.

There is one additional version that works like the Queue interface, where the data type is a native 'int' instead of interface{}. It uses the circular buffer approach but since its not using interface{} it might have better performance. The analysis will tell us that.

- NativeIntQueue
  - queue using a circular buffer with mutex/condition variable
  - [queue_native.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_native.go).
  - type safe for 'int'

At this point I'm not sure about the performance of any of these. We need to measure that.approaches.

### SynchronizedQueue

The Queue interface doesn't support some of the methods that are convenient for using in a threaded environment. The interface doesn't guarantee thread safety, nor does it specify blocking or non-blocking semantics. For that purpose, here is an extended interface that supports thread-safety and both blocking and non-blocking semantics. I provide an implementation of this interface that wraps a Queue with a SynchronizedQueue interface.

```go

// SynchronizedQueue is a queue with a bound on the number of elements in the queue
// Any implementation of this SHOULD promise thread-safety and the proper blocking semantics.
type  SynchronizedQueue interface {

	// add an element onto the tail queue
	// if the queue is full, the caller blocks
 	Put(value interface{})

	// add an element onto the tail queue
	// if the queue is full an error is returned
	TryPut(value interface{}) error

	// get an element from the head of the queue
	// if the queue is empty the caller blocks
	Get() interface{}

	// try to get an element from the head of the queue
	// if the queue is empty an error is returned
	TryGet() (interface{}, error)

	// current number of elements in the queue
    Len() int

	// capacity maximum number of elements the queue can hold
	Cap() int

	// close any resources (required for channel version)
	Close()

	// string representation
	fmt.Stringer
}
```

#### Synchronized Queue Using Channels

Of course, in the Go language, there are buffered channels, which literally are bounded queues. If you aren't familiar with Go channels, search for 'golang buffered channel' and there is lots of information. The official [Go Tour](https://tour.golang.org/concurrency/2) has a basic explanation. There are tons of references online about how to use buffered channels.

For channels there is an implementation of SynchronizedQueue without requiring the wrapper. There is no need for the Mutex/Cond support when you use channels.

See file [queue_channel.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_channel.go).

One other note about the channel version. If the queue is to be discarded at some point in execution but the program continues, the channel must be closed and all remaining data must be received so it doesn't leak. The interface has a Close() method that can be used by the application to close the channel when the queue is no longer needed. Only a producer calling Put may call Close(). A closed channel will panic if a Put is attempted. Once a channel is closed, any remaining data in the channel may still be accessible to Get's.

#### SynchronizedQueue Using Mutex/Condition Variable

There is an implementation of the SynchronizedQueue interface in the file **queue_sync.go**. This implementation takes a Queue and wraps it with the proper sync methods. In this case it uses the common Mutex/Condition variable approach. A function to create an instance of this interface is provided.

```go
// NewSynchronizedQueue is a factory for creating synchronized queues
// it takes an instance of the Queue interface and wraps it using
// a mutex and condition variable
// returns an instance of SynchronizedQueue
func NewSynchronizedQueue(q Queue) SynchronizedQueue
```

##### Mutex/Condition Variable

If you are familiar with the Mutex/Condition Variable paradigm, you can skip this section.

A Condition Variable is a synchronization object that allow threads to wait until a condition occurs. It does this in conjunction with a Mutex to provide mutual exclusion. The implementation of a interthread queue is a typical usage of a Mutex/Condition variable pair.

In the Go language, the standard library package 'sync' provides both _Mutex_ and _Cond_. It works like this:

```go

// declare a Mutex and A Cond
var mtx sync.Mutex
var cvr *sync.Cond

// initialize them
mtx = sync.Mutex{}
cvr = sync.NewCond(&mtx)
// the mutex is attached to the condition variable and after that
// it is accessed through the L.Lock() and L.Unlock() functions

// a function that responds if a condition is fulfilled
func f(cvr *sync.Cond) {
	// lock the mutex
	cvr.L.Lock()

	// make sure the mutex is unlocked when the function returns
	defer cvr.L.Unlock()

	// the mutex is locked at this point so its thread safe
	// to manipulate the data structure or resource

	// while some condition is not true, execute the loop
	for condition == false {
		// this function unlocks the associated mutex
		// and blocks the goroutine on the condition variable
		// It will wake up if some other goroutine calls Signal
		// on the condition variable
		cvr.Wait()
	}

	// the mutex is locked at this point so its thread safe
	// to manipulate the data structure or resource

	// the condition is now true, so do whatever is needed

	// signal a waiter if any
	// this wakes up one other goroutine that is blocked on the
	// condition variable. There is also Broadcast which will
	// wake up all goroutines waiting on the condition variable
	cvr.Signal()

	// the defer unlock releases the mutex
}
```

#### SynchronizedQueue Implementation

The file [queue_sync.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_sync.go) has the wrapper for a Queue that provide thread-safety using Mutex/Condition Variable support.

```go
// SynchronizedQueueImpl is an implementation of the SynchronizedQueue interface
// using a Mutex and 2 condition variables.
type SynchronizedQueueImpl struct {
	queue Queue	        // some data structure for backing the queue
	mtx sync.Mutex      // a mutex for mutual exclusion
	putcv *sync.Cond    // a condition variable for controlling Puts
	getcv *sync.Cond    // a condition variable for controlling Gets
}
```

Note that there are two condition variables and a single mutex. That is to support blocking on both ends, Put and Get. The single mutex protects the data structure in both the Put and Get calls. Put operations block on the **putcv** condition variable and signals the **getcv** condition variable when the Put adds an element to the queue. A Get operation works in the opposite direction, blocking on the **getcv** condition variable and signalling the **putcv** condition variable when an element is removed from the queue.

The non-blocking TryPut and TryGet operations still need to signal their opposite condition variable in case the other end is using the blocking version.

#### Queue Using a slice of interface{}

Hey, slices can act like queues. append to end, \[1:\] from front. In this case the slice is not preallocated. The
Put's and Get's modify the slice dynamically. Appending to the end is probably not too bad, but popping the front may be pretty ugly. We will see in the benchmark test.

See file [queue_slice.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_slice.go).

```go
// SliceQueue backed by a slice
// this version appends and removes elements so the slice grows and shrinks
// memory is not preallocated
// see CircularQueue for a version that preallocates a slice of capacity
type SliceQueue struct {
	slice []interface{}
	capacity int	     // maximum number of elements the queue can hold
}

// ...

// synchronized Slice queue
func NewSyncSlice(cap int) SynchronizedQueue {
	var sq Queue
	var bq SynchronizedQueue

	// create the queue
	sq = NewSliceQueue(cap)

	// wrap it with synchronization
	bq = NewSynchronizedQueue(sq)

	return bq
}

```

#### Queue Using container/list with interface{}

In this implementation the container/list data structure is used. In hindsight using a List is probably not the best approach since the queue is bounded so it doesn't need the flexibility of a List to shrink and grow. We'll see in the analysis.

See file [queue_list.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_list.go).

```go
// ListQueue
// the container/list data structure supports the semantics and methods
// needed for the Queue interface, with the exception of a capacity.
type ListQueue struct {
	list *list.List	 // contains the elements currently in the queue
	capacity int	 // maximum number of elements the queue can hold
}

// ...

// Synchronized List Queue Factory
func NewSyncList(cap int) SynchronizedQueue {
	var lq Queue
	var bq SynchronizedQueue

	// create a ListQueue
	lq = NewListQueue(cap)

	// wrap it with a SynchronizedQueue
	bq = NewSynchronizedQueue(lq)

	return bq
}
```

#### Queue Using container/ring with interface{}

In this implementation the container/ring data structure is used.

See file [queue_ring.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_ring.go).

```go
// RingQueue - a Queue backed by a container/ring
type RingQueue struct {
	ring *ring.Ring	 // preallocated ring for all slots in the queue
	head *ring.Ring  // head of the queue
	tail *ring.Ring  // tail of the queue
	capacity int     // maximum number of elements the ring can hold
	length int       // current number of element in the ring
}

// ...

// Synchronized RingQueue factory
func NewSyncRing(cap int) SynchronizedQueue {
	var rq Queue
	var bq SynchronizedQueue

	// create a ring queue
	rq = NewRingQueue(cap)

	// wrap it with a SynchronizedQueue
	bq = NewSynchronizedQueue(rq)

	return bq
}
```

#### Queue Using circular buffer with interface{}

This version uses a homegrown circular buffer as the queue data structure. Just guessing it should have better performance than the list. However it still uses interface{} for data elements so it might have some overhead for that vs a native data type.

See file [queue_circular.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_circular.go).

```go
// CircularQueue
type CircularQueue struct {
	queue []interface{}	// data
	head     int		// items are pulled from the head
	tail     int		// items are pushed to the tail
	length   int		// current number of elements in the queue
	capacity int	    // maximum allowed elements total
}

// ...

// Synchronized CircularQueue Factory
func NewSyncCircular(cap int) SynchronizedQueue {
	var cq Queue
	var bq SynchronizedQueue

	// create a circular queue
	cq = NewCircularQueue(cap)

	// wrap it with a synchronized queue
	bq = NewSynchronizedQueue(cq)

	return bq
}
```

#### Synchronized Queue using circular buffer with native ints

This version uses a circular buffer as the queue data structure. It is almost identical to the previous circular buffer version with the exception it only supports 'int' elements. I'm guessing that this may be a bit faster than the empty interface version. This version is not compatible with the SynchronizedQueue interface so it has its own mutual exclusion support.

See file [queue_native.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_native.go).

```go
// NativeIntQueue iis a type specific implementation
type NativeIntQueue struct {
	queue []int			// data
	head     int		// items are pulled from the head
	tail     int		// items are pushed to the tail
	length   int		// current number of elements in the queue
	capacity int	    // maximum allowed elements total
	mtx sync.Mutex      // a mutex for mutual exclusion
	putcv *sync.Cond    // a condition variable for controlling Puts
	getcv *sync.Cond    // a condition variable for controlling Gets
}

// NewNativeQueue is a factory for creating queues
// that use a condition variable and circular buffer
// for the specific type. In this case 'int'.
func NewNativeQueue(size int) *NativeIntQueue {
	var nvq NativeIntQueue

	// allocate the whole slice during init
	nvq.queue = make([]int,size,size)
	nvq.head = 0
	nvq.tail = 0
	nvq.length = 0
	nvq.capacity = size
	nvq.mtx = sync.Mutex{}
	nvq.putcv = sync.NewCond(&nvq.mtx)
	nvq.getcv = sync.NewCond(&nvq.mtx)

	return &nvq
}

```

#### Queue Using container/heap with interface{} (PriorityQueue)

In this implementation the container/heap structure is used. This is a special case. The implementation
creates a priority list. It is modeled after the [PriorityQueue example](https://golang.org/pkg/container/heap/#example__priorityQueue) from the standard library. So its not FIFO, instead the elements are ordered by priority. However it still complies with the Queue interface and can be used with the SynchronizedQueue wrapper.

This implementation requires a separate set of tests because the other ones use plain old ints for their data but this one requires a PriorityItem with both a **value interface{}** and **priority int**. It could be implemented with an int that represents both the value and priority but then its value won't be type agnostic.

file [queue_priority.go](https://github.com/dmh2000/go_sync_queue/blob/main/queue_priority.go) which contains the implementation of the Queue required by SynchronizedQueue.

```go
// one item in the priority queue
type PriorityItem struct {
	value interface{}
	priority int
}

// ...

// Synchronized PriorityQueue factory
func NewSyncPriorityQueue(cap int) SynchronizedQueue {
	var pq Queue
	var bq SynchronizedQueue

	// create the priority queue
	pq = NewPriorityQueue(cap)

	// wrap it in the synchronized queue
	bq = NewSynchronizedQueue(pq)

	return bq
}
```

### Testing

Three files have test code using the Go native test framework.

- run ./test.sh to execute all the tests

#### Synchronous Tests

The file [sync_test.go](https://github.com/dmh2000/go_sync_queue/blob/main/sync_test.go) contains tests where non-blocking TryPut's and TryGet's are performed in a single goroutine environment. These provide a test that the basic Try operations seem to work. They also test the capacity and length limits of the queues. There is no blocking in these tests.

As of this writing, all the synchronous tests pass.

```bash
# -v : verbose output
# -run <regex> : run only test functions match the regex
# . build with all files in current directory
$ go test -v -run Test.*Sync . | tee sync.out
=== RUN   TestPrioritySync
--- PASS: TestPrioritySync (0.00s)
=== RUN   TestChannelSync
--- PASS: TestChannelSync (0.00s)
=== RUN   TestQueueNativeSync
--- PASS: TestQueueNativeSync (0.00s)
=== RUN   TestCircularQueueSync
--- PASS: TestCircularQueueSync (0.00s)
=== RUN   TestListSync
--- PASS: TestListSync (0.00s)
=== RUN   TestRingSync
--- PASS: TestRingSync (0.00s)
=== RUN   TestSliceSync
--- PASS: TestSliceSync (0.00s)
=== RUN   TestStringsSync
    sync_test.go:201: ChannelQ Len:1 Cap:8
    sync_test.go:205: SynchronizedQueue:CircularQueue Len:1 Cap:8
    sync_test.go:209: SynchronizedQueue:ListQueue Len:1 Cap:8
    sync_test.go:213: SynchronizedQueue:RingQueue Len:1 Cap:8
    sync_test.go:217: SynchronizedQueue:PriorityQueue Len:1 Cap:8
    sync_test.go:221: SynchronizedQueue:SliceQueue Len:1 Cap:8
    sync_test.go:226: NativeIntQueue Len:1 Cap:8
--- PASS: TestStringsSync (0.00s)
PASS
ok      dmh2000.xyz/queue
```

#### Asynchronous Tests

The file [async_test.go](https://github.com/dmh2000/go_sync_queue/blob/main/async_test.go) contains tests where blocking Put's and Get's are performed in separate goroutines, one as a producer, one as a consumer. There are two versions of the tests, one where the Put and Get loops have no delays in them. The second is similar but with a random delay before each Put and Get. Intended to check that the blocking and wakeups are working properly.

As of this writing, all the synchronous tests pass.

```bash
# -v : verbose output
# -run <regex> : run only test functions match the regex
# . build with all files in current directory
$ go test -v -run Test.*Async .
=== RUN   TestChannelAsync
--- PASS: TestChannelAsync (0.27s)
=== RUN   TestListAsync
--- PASS: TestListAsync (0.24s)
=== RUN   TestCircularAsync
--- PASS: TestCircularAsync (0.25s)
=== RUN   TestRingAsync
--- PASS: TestRingAsync (0.23s)
=== RUN   TestSliceAsync
--- PASS: TestSliceAsync (0.25s)
=== RUN   TestComboAsync
--- PASS: TestComboAsync (0.21s)
=== RUN   TestNativeAsync
--- PASS: TestNativeAsync (0.22s)
PASS
ok  	dmh2000.xyz/queue	
```

#### Benchmarks

The file [benchmark_test.go](https://github.com/dmh2000/go_sync_queue/blob/main/benchmark_test.go) contains tests both synchronous and asynchronous version similar to the Async tests, but intended to be used as benchmarks for timing and memory usage.

```bash
# -v : verbose output
# -run <regex> : run only test functions match the regex
# -bench : run benchmarks
# . : build with all files in current directory
# -benchmem : measure memory usage
# -memprofile <file> : write memory usage information to <file>
# -cpuprofile <file> : write cpu usage information to <file>
$ go test -v -run Benchmark.* -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out
goos: linux
goarch: amd64
pkg: dmh2000.xyz/queue
cpu: AMD Ryzen 7 3700X 8-Core Processor             
BenchmarkQueueChannelSync
BenchmarkQueueChannelSync-16     	 1330677	       915.3 ns/op	     424 B/op	       3 allocs/op
BenchmarkListSync
BenchmarkListSync-16             	  539108	      2115 ns/op	    1200 B/op	      25 allocs/op
BenchmarkCircularSync
BenchmarkCircularSync-16         	  920787	      1260 ns/op	     560 B/op	       5 allocs/op
BenchmarkRingSync
BenchmarkRingSync-16             	  597390	      1842 ns/op	     864 B/op	      24 allocs/op
BenchmarkSliceSync
BenchmarkSliceSync-16            	  663648	      1615 ns/op	    1216 B/op	      10 allocs/op
BenchmarkQueueNativeSync
BenchmarkQueueNativeSync-16      	 1470663	       812.7 ns/op	     368 B/op	       4 allocs/op
BenchmarkQueueChannelAsync
BenchmarkQueueChannelAsync-16    	  458401	      2386 ns/op	     520 B/op	       6 allocs/op
BenchmarkListAsync
BenchmarkListAsync-16            	  319388	      3505 ns/op	    1296 B/op	      28 allocs/op
BenchmarkCircularAsync
BenchmarkCircularAsync-16        	  426242	      2646 ns/op	     656 B/op	       8 allocs/op
BenchmarkRingAsync
BenchmarkRingAsync-16            	  304796	      3514 ns/op	     960 B/op	      27 allocs/op
BenchmarkSliceAsync
BenchmarkSliceAsync-16           	  372702	      3019 ns/op	    1312 B/op	      13 allocs/op
BenchmarkQueueNativeAsync
BenchmarkQueueNativeAsync-16     	  486388	      2223 ns/op	     440 B/op	       7 allocs/op
PASS
ok  	dmh2000.xyz/queue	15.664s
```

Here are the synchronous tests sorted by time/iter (second column). Least number indicates slowest. The number of iterations is skewed by the cost of preallocating in some cases so that isn't the best measure

<pre>
BenchmarkQueueNativeSync-4     628731   1846.00 ns/op      368 B/op    4 allocs/op
BenchmarkQueueChannelSync-4    619838   2016.00 ns/op      424 B/op    3 allocs/op
BenchmarkCircularSync-4        443493   2750.00 ns/op      560 B/op    5 allocs/op
BenchmarkRingSync-4            459682   3233.00 ns/op      864 B/op   24 allocs/op
BenchmarkSliceSync-4           421218   3412.00 ns/op     1216 B/op   10 allocs/op
BenchmarkListSync-4            362503   3680.00 ns/op     1200 B/op   25 allocs/op
</pre>

- an 'op' is one execution of one iteration of the benchmark test in the **for i:=0;i<b.N;i++** loop
- Looks like channels and native type versions work best, kind of as expected.
- The versions using containers suffered from being more general purpose.
- It also appears that fewer bytes allocated per op and number of allocations may be significant.

#### Race Detector

The file [race_test.go](https://github.com/dmh2000/go_sync_queue/blob/main/race_test.go) is uses the same test routines as the async_test.go but with a much larger queue size. It is used to run for a longer time to allow the go race detector to see more possible problems. Since one of the async tests uses time delays, this test will run for a few minutes.

```bash
$  go test -v -race *.go 
=== RUN   TestChannelAsync
--- PASS: TestChannelAsync (0.30s)
=== RUN   TestListAsync
--- PASS: TestListAsync (0.20s)
=== RUN   TestCircularAsync
--- PASS: TestCircularAsync (0.19s)
=== RUN   TestRingAsync
--- PASS: TestRingAsync (0.28s)
=== RUN   TestSliceAsync
--- PASS: TestSliceAsync (0.21s)
=== RUN   TestComboAsync
--- PASS: TestComboAsync (0.21s)
=== RUN   TestNativeAsync
--- PASS: TestNativeAsync (0.25s)
=== RUN   TestPrioritySync
--- PASS: TestPrioritySync (0.00s)
=== RUN   TestChannelRace
--- PASS: TestChannelRace (2.82s)
=== RUN   TestListRace
--- PASS: TestListRace (2.58s)
=== RUN   TestCircularRace
--- PASS: TestCircularRace (2.63s)
=== RUN   TestRingRace
--- PASS: TestRingRace (2.61s)
=== RUN   TestSliceRace
--- PASS: TestSliceRace (2.72s)
=== RUN   TestComboRace
--- PASS: TestComboRace (2.66s)
=== RUN   TestNativeRace
--- PASS: TestNativeRace (2.57s)
=== RUN   TestChannelSync
--- PASS: TestChannelSync (0.00s)
=== RUN   TestQueueNativeSync
--- PASS: TestQueueNativeSync (0.00s)
=== RUN   TestCircularQueueSync
--- PASS: TestCircularQueueSync (0.00s)
=== RUN   TestListSync
--- PASS: TestListSync (0.00s)
=== RUN   TestRingSync
--- PASS: TestRingSync (0.00s)
=== RUN   TestSliceSync
--- PASS: TestSliceSync (0.00s)
=== RUN   TestStringsSync
    sync_test.go:201: ChannelQ Len:1 Cap:8
    sync_test.go:205: SynchronizedQueue:CircularQueue Len:1 Cap:8
    sync_test.go:209: SynchronizedQueue:ListQueue Len:1 Cap:8
    sync_test.go:213: SynchronizedQueue:RingQueue Len:1 Cap:8
    sync_test.go:217: SynchronizedQueue:PriorityQueue Len:1 Cap:8
    sync_test.go:221: SynchronizedQueue:SliceQueue Len:1 Cap:8
    sync_test.go:226: NativeIntQueue Len:1 Cap:8
--- PASS: TestStringsSync (0.00s)
PASS
ok  	command-line-arguments	21.239s
```

No race conditions were detected. Since this is an artificial test it might not find problems that would occur in an actual application with a different call sequence. Mutual exclusion is tricky! Don't assume an implementation works because it just looks right.

#### Analysis

##### Memory

After running the benchmark tests above, the file mem.out contains information about the memory usage of the benchmarks. This file can be analyzed using the pprof go tool.

```bash
# run the pprof tool with input file mem.out
# (pprof) top10 : shows top 10 locations for memory allocation
$ go tool pprof -top  -nodecount=15 mem.out
File: queue.test
Type: alloc_space
Time: Feb 10, 2024 at 1:32am (UTC)
Showing nodes accounting for 5.98GB, 99.94% of 5.98GB total
Dropped 29 nodes (cum <= 0.03GB)
Showing top 15 nodes out of 37
      flat  flat%   sum%        cum   cum%
    1.11GB 18.49% 18.49%     1.11GB 18.49%  dmh2000.xyz/queue.NewChannelQueue
    1.01GB 16.88% 35.37%     1.01GB 16.88%  dmh2000.xyz/queue.(*SliceQueue).Push
    0.89GB 14.80% 50.17%     0.89GB 14.80%  sync.NewCond (inline)
    0.81GB 13.50% 63.68%     0.81GB 13.50%  container/list.(*List).insertValue (inline)
    0.66GB 11.02% 74.70%     1.04GB 17.31%  dmh2000.xyz/queue.NewNativeQueue
    0.54GB  8.98% 83.68%     0.54GB  8.98%  container/ring.New (inline)
    0.47GB  7.87% 91.54%     0.47GB  7.87%  dmh2000.xyz/queue.NewCircularQueue
    0.18GB  2.95% 94.49%     0.69GB 11.46%  dmh2000.xyz/queue.NewSynchronizedQueue
    0.16GB  2.71% 97.20%     0.16GB  2.71%  dmh2000.xyz/queue.asyncb1
    0.04GB  0.72% 97.92%     0.58GB  9.69%  dmh2000.xyz/queue.NewRingQueue
    0.04GB  0.65% 98.57%     0.04GB  0.65%  dmh2000.xyz/queue.asyncb2
    0.04GB   0.6% 99.17%     0.04GB   0.6%  container/list.New (inline)
    0.04GB  0.59% 99.76%     0.04GB  0.59%  dmh2000.xyz/queue.NewSliceQueue
    0.01GB  0.18% 99.94%     0.05GB  0.78%  dmh2000.xyz/queue.NewListQueue
         0     0% 99.94%     0.81GB 13.50%  container/list.(*List).PushBack
```

- Notes

  1. The slice and list versions both allocate a lot of memory and they do it per Put

  - The list isn't preallocated
  - There is probably an allocation for every Push (append) and insertValue statement
  - this hits the runtime after initialization.

  2.  The creation of the condition variable uses a lot of time

  - this is the combined usage over all the versions that use Mutex/Cond
  - see #5

  3. The creation of the container/ring in the Ring version is expensive

  - the Ring is preallocated
  - this hit is only during initialization of the Ring

  4. The Channel, Native, and Circular versions all preallocate during initialization
  5. Creating the SynchronizedQueue wrapper does some work

  - is appears that creating a Cond is a bit expensive
  - there is at least one system call involved there

  6. The slice queue does an allocation when it appends an element

  7. The rest of the hits are all less than 1%

- Conclusion
  - Don't use a container/list !

##### CPU

```bash
# run the pprof tool with input file cpu.out
# (pprof) top10 : shows top 10 locations for cpu usage
$ go tool pprof -top -nodecount=20  cpu.out
File: queue.test
Type: cpu
Time: Feb 10, 2024 at 1:31am (UTC)
Duration: 15.46s, Total samples = 13090ms (84.67%)
Showing nodes accounting for 8250ms, 63.03% of 13090ms total
Dropped 184 nodes (cum <= 65.45ms)
Showing top 20 nodes out of 156
      flat  flat%   sum%        cum   cum%
    1480ms 11.31% 11.31%     1480ms 11.31%  runtime.futex
     850ms  6.49% 17.80%     2790ms 21.31%  runtime.mallocgc
     700ms  5.35% 23.15%     2230ms 17.04%  dmh2000.xyz/queue.(*SynchronizedQueueImpl).Put
     610ms  4.66% 27.81%      610ms  4.66%  sync.(*Mutex).Lock
     520ms  3.97% 31.78%      520ms  3.97%  sync.(*Mutex).Unlock
     460ms  3.51% 35.29%      470ms  3.59%  runtime.unlock2
     440ms  3.36% 38.66%      980ms  7.49%  dmh2000.xyz/queue.(*SynchronizedQueueImpl).Get
     380ms  2.90% 41.56%      800ms  6.11%  dmh2000.xyz/queue.(*NativeIntQueue).Put
     360ms  2.75% 44.31%      360ms  2.75%  runtime.nextFreeFast (inline)
     260ms  1.99% 46.29%      650ms  4.97%  dmh2000.xyz/queue.(*NativeIntQueue).Get
     260ms  1.99% 48.28%     3940ms 30.10%  dmh2000.xyz/queue.b1
     260ms  1.99% 50.27%      260ms  1.99%  runtime.usleep
     260ms  1.99% 52.25%      260ms  1.99%  runtime.writeHeapBits.flush
     250ms  1.91% 54.16%      310ms  2.37%  runtime.lock2
     230ms  1.76% 55.92%      230ms  1.76%  dmh2000.xyz/queue.(*SynchronizedQueueImpl).Cap
     220ms  1.68% 57.60%      220ms  1.68%  runtime.memclrNoHeapPointers
     190ms  1.45% 59.05%      190ms  1.45%  sync.(*copyChecker).check (inline)
     180ms  1.38% 60.43%      630ms  4.81%  runtime.heapBitsSetType
     170ms  1.30% 61.73%      850ms  6.49%  runtime.chanrecv
     170ms  1.30% 63.03%      170ms  1.30%  runtime.releasem (inline)
```

- Notes
  - The preallocations in the New functions don't seem to hit the runtime
  - Most of the work is being done in system calls relating to mutual exclusion
  - so Mutexes and Conds are expensive to use
  - create a non-blocking version? Too hard for me
  - The Gets and Puts also hit the Mutex/Conds
  - the 'Try' functions are much cheaper
