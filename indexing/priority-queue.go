package indexing

type Record struct {
	workerId int    // ATTENTION: this is not index!
	workload int    // the approximated amount of data held by this worker
}

// PriorityQueue a minimum heap with fixed size.
//
// Please use NewPriorityQueue() as the safe constructor.
type PriorityQueue struct {
	data []*Record
	size int
}

func NewPriorityQueue(size int) *PriorityQueue {
	return &PriorityQueue{
		data: make([]*Record, size),
		size: size,
	}
}

// ExtractMin extracts the minimum element and returns it as a Record object.
func (pq *PriorityQueue) ExtractMin() *Record {
	//todo
	return nil
}

// Add adds a new element into the priority queue.
func (pq *PriorityQueue) Add(key, value int) bool {
	//todo
	return true
}
