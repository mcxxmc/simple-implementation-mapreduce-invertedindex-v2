package indexing

// Manager the manager that allocates tasks between all workers / collectors.
//
// Note that it does not hold direct pointers to the workers; it holds pointers to the communication channels (one for
// each worker) instead, to simulate the communication between different machines over network.
//
// Please use NewManager() as the safe constructor.
type Manager struct {
	workers      []chan *Msg    // the communication channels (one for each worker)
	listen       chan *Msg      // its own channel listening from workers
	alive        map[int]bool   // mark each worker as alive or dead; in the form of: workerId: bool
	idle         *Queue         // mark each worker as idle or busy
	workloads    *PriorityQueue // a minimum heap, which should return the worker with the slightest amounts of data.
	numOfWorkers int            // num of workers supervised
	id           int            // the id of the manager
}

// NewManager returns a pointer to a new Manager object, with default id 0.
func NewManager(numberOfWorkers int) *Manager {
	workers := make([]chan *Msg, numberOfWorkers)
	for i := range workers {
		workers[i] = make(chan *Msg, CollectorChanCapacity)
	}
	return &Manager{
		workers:      workers,
		listen:       make(chan *Msg, ManagerChanCapacity),
		alive:        make(map[int]bool),
		idle:         NewQueue(numberOfWorkers),
		workloads:    NewPriorityQueue(numberOfWorkers),
		numOfWorkers: numberOfWorkers,
		id:           0,
	}
}

// converts index to worker id.
func (manager *Manager) getWorkerId(index int) int {
	return index + 1
}

// Run runs the manager.
func (manager *Manager) Run() {
	// creates workers; by default, worker id is index + 1.
	// the worker indexes are pushed into and popped out from manager.idle in a FIFO fashion.
	for index, channel := range manager.workers {
		workerId := manager.getWorkerId(index)
		collector := NewCollector(channel, manager.listen, workerId)
		manager.alive[workerId] = true
		manager.idle.Push(workerId)
		go collector.Run()
	}

	//todo
}

// returns the worker id of an alive and idle worker for a new task.
//
// If no worker is idle right now, will return 0 and false.
func (manager *Manager) nextWorkerId() (int, bool) {
	var workerId int
	var exist bool
	for {
		workerId, exist = manager.idle.Pop()
		if !exist { // the queue is empty
			return 0, false
		}
		if manager.alive[workerId] {
			return workerId, true
		}
	}
}
