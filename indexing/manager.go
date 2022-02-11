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
	online       map[int]bool   // mark each worker as online or offline; in the form of: workerId: bool
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
		online:       make(map[int]bool),
		workloads:    NewPriorityQueue(numberOfWorkers),
		numOfWorkers: numberOfWorkers,
		id:           0,
	}
}

// converts index to worker id.
func (manager *Manager) getWorkerId(index int) int {
	return index + 1
}

// Initialize initializes the manager with the tasks (paths of the files to be processed).
func (manager *Manager) Initialize(tasks ...string) {
	//todo
}

// Run runs the manager.
//
// Must call Initialize() before this method.
func (manager *Manager) Run() {
	// creates workers; by default, worker id is index + 1.
	// the worker indexes are pushed into and popped out from manager.idle in a FIFO fashion.
	for index, channel := range manager.workers {
		workerId := manager.getWorkerId(index)
		collector := NewCollector(channel, manager.listen, workerId)
		manager.online[workerId] = true
		go collector.Run()
	}

	//todo
}

// returns the worker id of an online worker (with the slightest workload) for a new task.
//
// If no worker is online right now, will return 0 and false.
func (manager *Manager) nextWorkerId() (int, bool) {
	//todo
	return 0, false
}
