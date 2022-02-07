package indexing

// Manager the manager that allocates tasks between all workers / collectors.
//
// Note that it does not hold direct pointers to the workers; it holds pointers to the communication channels (one for
// each worker) instead, to simulate the communication between different machines over network.
//
// Please use NewManager() as the safe constructor.
type Manager struct {
	workers []chan *Msg  // the communication channels (one for each worker)
	listen chan *Msg     // its own channel listening from workers
	alive []bool 		 // mark each worker as alive or dead
	idle *Queue  		 // mark each worker as idle or busy
	roundRobinIndex int  // the next index in a round-robin fashion
	numOfWorkers int  	 // num of workers supervised
	id int               // the id of the manager
}

// NewManager returns a pointer to a new Manager object, with default id 0.
func NewManager(numberOfWorkers int) *Manager {
	workers := make([]chan *Msg, numberOfWorkers)
	for i := range workers {
		workers[i] = make(chan *Msg, CollectorChanCapacity)
	}
	return &Manager{
		workers: workers,
		listen: make(chan *Msg, ManagerChanCapacity),
		alive: make([]bool, numberOfWorkers),
		idle: NewQueue(numberOfWorkers),
		roundRobinIndex: 0,
		numOfWorkers: numberOfWorkers,
		id: 0,
	}
}

// converts index to worker id.
func (manager *Manager) getWorkerId(index int) int {
	return index + 1
}

// converts worker id to index.
func (manager *Manager) getIndex(workerId int) int {
	return workerId - 1
}

// Run runs the manager.
func (manager *Manager) Run() {
	// creates workers; by default, worker id is index + 1.
	// the worker indexes are pushed into and popped out from manager.idle in a FIFO fashion.
	for i, channel := range manager.workers {
		collector := NewCollector(channel, manager.listen, manager.getWorkerId(i))
		manager.alive[i] = true
		manager.idle.Push(i)
		go collector.Run()
	}

	//todo
}

// selects an alive worker for a new task. If no worker is idle right now, use the round-robin index.
func (manager *Manager) selectIndex() int {
	var r int
	var b bool
	for {
		r, b = manager.idle.Pop()
		if !b {  // the queue is empty
			index := manager.roundRobinIndex
			if index == manager.numOfWorkers - 1 {
				manager.roundRobinIndex = 0
			} else {
				manager.roundRobinIndex ++
			}
			return index
		}
		if manager.alive[r] {
			return r
		}
	}
}
