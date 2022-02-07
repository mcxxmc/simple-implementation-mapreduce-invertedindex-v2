package indexing

// Queue a first-in-first-out queue for int values.
type Queue struct {
	data []int
	size int    	// the max size
	numOfItems int  // the current number of items in the queue
	start int  	    // the left (lower) boundary, points to the index of the first item
}

// NewQueue returns a new queue object.
func NewQueue(size int) *Queue {
	return &Queue{data: make([]int, size), size: size, numOfItems: 0, start: 0}
}

// Capacity returns the max size of the queue
func (q *Queue) Capacity() int {
	return q.size
}

// NumberOfItems returns the current number of items in the queue
func (q *Queue) NumberOfItems() int {
	return q.numOfItems
}

// Push pushes a new item into the queue.
func (q *Queue) Push(val int) bool {
	if q.size == q.numOfItems {
		return false
	}
	index := q.start + q.numOfItems
	if index >= q.size {
		index -= q.size
	}
	q.data[index] = val
	q.numOfItems ++
	return true
}

// Pop pops an item out from the queue
func (q *Queue) Pop() (int, bool) {
	if q.numOfItems == 0 {
		return 0, false
	}
	r := q.data[q.start]
	if q.start == q.size - 1 {
		q.start = 0
	} else {
		q.start ++
	}
	q.numOfItems --
	return r, true
}

// Peak gets (instead of popping out) the first item in the queue.
func (q *Queue) Peak() (int, bool) {
	if q.numOfItems == 0 {
		return 0, false
	}
	return q.data[q.start], true
}
