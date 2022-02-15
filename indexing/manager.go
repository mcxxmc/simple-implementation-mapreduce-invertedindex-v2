package indexing

import (
	"errors"
	"fmt"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"strconv"
	"time"
)

// Manager the manager that allocates tasks between all workers / collectors.
//
// Note that it does not hold direct pointers to the workers; it holds pointers to the communication channels (one for
// each worker) instead, to simulate the communication between different machines over network.
//
// Please use NewManager() as the safe constructor.
type Manager struct {
	workers      	map[int]chan *Msg   // the communication channels (one for each worker)
	online       	map[int]bool   		// mark each worker as online or offline; in the form of: workerId: bool
	workloadRecords map[int]int			// workerId: workload
	listen       	chan *Msg      		// its own channel listening from workers
	workloads    	*PriorityQueue 		// a minimum heap, which should return the worker with the slightest amounts of data.
	numOfWorkers 	int            		// current number of workers supervised
	id           	int            		// the id of the manager
}

// NewManager returns a pointer to a new Manager object, with default id 0.
func NewManager(numberOfWorkers int) *Manager {
	return &Manager{
		workers:      make(map[int]chan *Msg),
		online:       make(map[int]bool),
		workloadRecords: make(map[int]int),
		listen:       make(chan *Msg, ManagerChanCapacity),
		workloads:    NewPriorityQueue(numberOfWorkers),
		numOfWorkers: numberOfWorkers,
		id:           0,
	}
}

// Run runs the manager with the tasks (paths of the files to be processed).
func (manager *Manager) Run(tasks []string) {
	// creates workers; by default, worker id starts from 1.
	for workerId := 1; workerId <= manager.numOfWorkers; workerId ++  {
		collector, err := manager.newWorker(workerId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go collector.Run()
	}

	manager.releaseWorkFirstPatch(tasks)

	// for reallocating, cleaning up and retiring
	deliveryToRequest := manager.numOfWorkers / 2
	deliveryRequested := 0
	deliveryReceived := 0

	running := true
	for running {
		select {
		case msg := <- manager.listen:
			switch msg.Typ {
			case MsgCollectorIOCompleted:
				oldWorkerId := msg.Id
				fmt.Println("all tasks completed, final step by worker ", oldWorkerId)
				manager.retireWorker(oldWorkerId)
				running = false

			case MsgCollectorDelivery:
				oldWorkerId := msg.Id
				fmt.Println("manager receives data from worker ", oldWorkerId)
				manager.retireWorker(oldWorkerId)
				manager.reallocate(oldWorkerId, msg.Data.(common.NativeRecords))
				deliveryReceived ++

			default:
				fmt.Println("manager receives unknown data type")
			}

		case <- time.After(time.Second):	// remove the time.After() here to allow the manager to do checking whenever it is idle
			fmt.Println("manager idle, moving to the next stage; current online workers: ", manager.numOfWorkers)

			if deliveryRequested == deliveryToRequest {
				if deliveryReceived < deliveryRequested {
					fmt.Println("not all requested data are delivered yet; waiting for the next loop")
					continue
				} else {
					deliveryToRequest = manager.numOfWorkers / 2
					deliveryRequested = 0
					deliveryReceived = 0
				}
			}

			switch manager.numOfWorkers {
			case 0:
				continue

			case 1:
				manager.requestSortSave2Disk()

			default:  // > 1
				manager.requestData()
				deliveryRequested ++
			}
		}
	}

	return
}

// creates and returns the pointer to a new worker
func (manager *Manager) newWorker(workerId int) (*Collector, error) {
	if _, exist := manager.workers[workerId]; exist {
		return nil, errors.New("workerId already exist: " + strconv.Itoa(workerId))
	}
	channel := make(chan *Msg, CollectorChanCapacity)
	collector := NewCollector(channel, manager.listen, workerId)
	manager.workers[workerId] = channel
	manager.online[workerId] = true
	manager.workloads.Add(workerId, 0)
	return collector, nil
}

// releases the first patch of work
func (manager *Manager) releaseWorkFirstPatch(tasks []string) {
	for _, task := range tasks {
		record := manager.nextWorker()
		workload := getWorkload(task) + record.Workload
		manager.workers[record.WorkerId] <- NewMsgCountFreq(task, manager.id)
		fmt.Println("task ", task, " assigned to worker ", record.WorkerId, ", current workload in KB: ", workload)
		manager.workloads.Add(record.WorkerId, workload)
		manager.workloadRecords[record.WorkerId] = workload
	}
}

// retires a worker by id
func (manager *Manager) retireWorker(oldWorkerId int) {
	manager.workers[oldWorkerId] <- NewMsgClearData(manager.id)
	manager.workers[oldWorkerId] <- NewMsgDismissWorker(manager.id)
	manager.online[oldWorkerId] = false
}

// reallocates the data from a worker to be retired to a still online one
func (manager *Manager) reallocate(prevWorker int, data common.NativeRecords) {
	record := manager.nextWorker()
	workload := manager.workloadRecords[prevWorker] + record.Workload
	manager.workers[record.WorkerId] <- NewMsgCombineFreq(data, prevWorker)
	fmt.Println("manager redirects data from worker ", prevWorker, " to worker ", record.WorkerId)
	manager.workloads.Add(record.WorkerId, workload)
	manager.workloadRecords[record.WorkerId] = workload
}

// requests data from a worker with the least workload
func (manager *Manager) requestData() {
	record := manager.nextWorker()
	fmt.Println("manager requests data from worker ", record.WorkerId)
	manager.workers[record.WorkerId] <- NewMsgDeliverData(manager.id)
	manager.numOfWorkers --
}

// requests the final worker to sort and save the data onto the disk
func (manager *Manager) requestSortSave2Disk() {
	record := manager.nextWorker()
	fmt.Println("manager sends final task to worker ", record.WorkerId)
	manager.workers[record.WorkerId] <- NewMsgSortSave2Disk(common.InvertedIndexSavePath, manager.id)
	manager.numOfWorkers --
}

// returns the Record of an online worker (with the slightest workload) for a new task.
//
// If no worker is online right now, will cause panic.
func (manager *Manager) nextWorker() *Record {
	record := manager.workloads.ExtractMin()
	for record != nil && !manager.online[record.WorkerId] {
		record = manager.workloads.ExtractMin()
	}
	if record == nil {
		panic("No worker online")
	}
	return record
}


// returns the number of KBs of a file.
func getWorkload(file string) int {
	f, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
	}
	return int(f.Size() / 1024)
}
