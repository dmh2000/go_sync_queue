package queue

// this is the priority queue example from container/heap
// it is used as-is
// https://golang.org/pkg/container/heap/#example__priorityQueue
// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	index int // the index of item in the heap
}

// A PrioritySliceX implements heap.Interface and holds Items.
type PrioritySliceX []*Item

func (pq PrioritySliceX) Len() int { return len(pq) }

func (pq PrioritySliceX) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PrioritySliceX) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PrioritySliceX) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PrioritySliceX) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// the update method is left out for this implementation