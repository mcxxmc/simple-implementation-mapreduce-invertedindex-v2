package indexing

import "fmt"

type Record struct {
	WorkerId int // ATTENTION: this is not index!
	Workload int // the approximated amount of data held by this worker
}

// PriorityQueue a minimum heap with fixed size.
//
// Please use NewPriorityQueue() as the safe constructor.
type PriorityQueue struct {
	data 	[]*Record
	limit 	int			// the max number of elements
	num  	int			// the current number of elements
}

func NewPriorityQueue(size int) *PriorityQueue {
	return &PriorityQueue{
		data: make([]*Record, size + 1),
		limit: size,
		num:  0,
	}
}

// Num returns the num of elements in the priority queue.
func (pq *PriorityQueue) Num() int {
	return pq.num
}

// ExtractMin extracts the minimum element and returns it as a Record object.
//
// Will return nil if this is no such element.
func (pq *PriorityQueue) ExtractMin() *Record {
	if pq.num < 1 {
		return nil
	}
	r := pq.data[1]
	pq.data[1] = pq.data[pq.num]
	pq.num--
	pq.minHeapify(1)
	return r
}

// Add adds a new element into the priority queue.
func (pq *PriorityQueue) Add(key, value int) bool {
	record := &Record{WorkerId: key, Workload: value}
	pq.num ++
	if pq.num > pq.limit {
		fmt.Println("priority full")
		pq.num --
		return false
	}
	pq.reorder(pq.num, record)
	return true
}

func (pq *PriorityQueue) swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}

// a method used for "sorting" any existing subtree of the heap from the root.
//
// i: the index of the current node
func (pq *PriorityQueue) minHeapify(i int) {
	left := i << 1
	right := left + 1
	smallest := i
	if left <= pq.num && pq.data[smallest].Workload > pq.data[left].Workload {
		smallest = left
	}
	if right <= pq.num && pq.data[smallest].Workload > pq.data[right].Workload {
		smallest = right
	}
	if smallest != i {
		pq.swap(i, smallest)
		pq.minHeapify(smallest)
	}
}

// inserts a new record to the end of the heap and sort the subtree from the leaf.
func (pq *PriorityQueue) reorder(index int, record *Record) {
	pq.data[index] = record
	for index > 1 && pq.data[index].Workload < pq.data[index / 2].Workload {
		pq.swap(index, index / 2)
		index = index / 2
	}
}
