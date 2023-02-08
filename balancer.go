package goBalancer

import "sync"

type runnable interface {
	// Run() is the function that will be executed by the load balancer.
	Run()
}

// one item in the queue
type queueItem[T runnable] struct {
	lck       *sync.Mutex
	id        int
	item      T
	isRunning bool
}

type loadBalancer[T runnable] struct {
	// Makes sure that the job ID is unique
	ids int
	// Lock for accessing ids variable
	idlck *sync.Mutex
	// Queue of jobs
	l []*queueItem[T]
	// Lock for accessing queue l
	// qlck should only be acquired by functions implemented for queue
	qlck *sync.Mutex
	// Maximal number of jobs to execute at the same time
	executeNjobs int
	// Total number of jobs running now
	totalRunningJobs int
	// Lock for accessing totalRunningJobs varialbe
	statslck *sync.Mutex
	// WaitGroup for waiting for all jobs to finish
	universalWg *sync.WaitGroup
}

// Config is a struct that holds the initial configurations for the load balancer.
type Config[T any] struct {
	// Execute is a function that gets executed on the queue.
	// Execute func(v *T) error
	// Load Balancer will Execute at max N jobs at the same time.
	MaxJobsToExecuteAtSameTime int
}

// NewGoBalancer() returns a new load balancer instance.
func NewGoBalancer[T runnable](c Config[T]) *loadBalancer[T] {

	if c.MaxJobsToExecuteAtSameTime <= 0 {
		c.MaxJobsToExecuteAtSameTime = 1
	}

	return &loadBalancer[T]{
		executeNjobs: c.MaxJobsToExecuteAtSameTime,
		ids:          0,
		idlck:        &sync.Mutex{},
		l:            []*queueItem[T]{},
		qlck:         &sync.Mutex{},
		statslck:     &sync.Mutex{},
		universalWg:  &sync.WaitGroup{},
	}
}

// AddJob() adds a new job to the load balancer queue. If the number of jobs running is less than the maximum number
// of jobs to execute at the same time, the job will be executed immediately. Otherwise, the job will be added to the queue.
// And will be executed when a slot gets freed.
// The function returns the job ID.
func (lb *loadBalancer[T]) AddJob(item T) int {
	jobID := lb.newJobID()

	lb.addToQueue(queueItem[T]{
		lck:       &sync.Mutex{},
		id:        jobID,
		item:      item,
		isRunning: false})

	if lb.hasFreeSlot() {
		// execute job
		lb.execute(jobID)
	}

	return jobID
}

// RemoveJobByID() removes a job from the queue by its ID.
// The function returns the removed job.
func (lb *loadBalancer[T]) RemoveJobByID(ID int) *T {
	return lb.removeFromQueue(ID)
}

// TotalJobs() returns the total number of jobs.
func (lb *loadBalancer[T]) TotalJobs() int {
	lb.qlck.Lock()
	defer lb.qlck.Unlock()
	return len(lb.l)
}

// TotalJobs() returns the total number of running jobs.
func (lb *loadBalancer[T]) TotalRunningJobs() int {
	lb.qlck.Lock()
	defer lb.qlck.Unlock()
	return len(lb.l)
}

func (lb *loadBalancer[T]) WaitForRunningJobs() {
	lb.universalWg.Wait()
}

// newJobID() returns a new unique job ID.
func (lb *loadBalancer[T]) newJobID() int {
	lb.idlck.Lock()
	defer lb.idlck.Unlock()

	lb.ids++
	return lb.ids
}

func (lb *loadBalancer[T]) incTotalRunningJobs() {
	lb.statslck.Lock()
	lb.totalRunningJobs++
	lb.statslck.Unlock()
}

func (lb *loadBalancer[T]) decTotalRunningJobs() {
	lb.statslck.Lock()
	lb.totalRunningJobs--
	lb.statslck.Unlock()
}

func (lb *loadBalancer[T]) triggerScheduler() {

	if lb.hasFreeSlot() {
		id := lb.nextScheduledJobID()
		// All jobs are running
		if id <= 0 {
			return
		}

		lb.execute(id)
	}

}

func (lb *loadBalancer[T]) nextScheduledJobID() int {

	for _, v := range lb.l {
		if !v.isRunning {
			return v.id
		}
	}

	return -1
}

// execute() executes the job.
func (lb *loadBalancer[T]) execute(jobID int) {
	lb.universalWg.Add(1)
	defer lb.universalWg.Done()

	lb.incTotalRunningJobs()
	defer lb.decTotalRunningJobs()

	var index int = -1
	for i, v := range lb.l {
		if v.id == jobID {
			if v.isRunning {
				return
			} else {
				index = i
				break
			}
		}
	}

	if index == -1 {
		return
	}

	lb.l[index].lck.Lock()

	lb.l[index].isRunning = true
	lb.l[index].item.Run()
	// After job is done, remove it from the queue and trigger the scheduler.
	lb.l[index].isRunning = false
	lb.l[index].lck.Unlock()

	lb.removeFromQueue(jobID)

	lb.triggerScheduler()
}

func (lb *loadBalancer[T]) addToQueue(item queueItem[T]) {
	lb.qlck.Lock()
	lb.l = append(lb.l, &item)
	lb.qlck.Unlock()
}

func (lb *loadBalancer[T]) removeFromQueue(ID int) *T {
	for i, v := range lb.l {
		if v.id == ID {
			lb.qlck.Lock()
			lb.l = append(lb.l[:i], lb.l[i+1:]...)
			lb.qlck.Unlock()
			return &v.item
		}
	}
	return nil
}

func (lb *loadBalancer[T]) hasFreeSlot() bool {
	lb.statslck.Lock()
	defer lb.statslck.Unlock()
	return lb.totalRunningJobs < lb.executeNjobs
}

func toPtr[T any](v T) *T {
	return &v
}
