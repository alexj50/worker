package worker

import (
	"sync"
	"testing"
)

type Job interface {
	perform()
	performTest(t *testing.T)
}

type Worker struct {
	MaxQueue, MaxWorkers int
	Testing *testing.T
}

// A buffered channel that we can send work requests on.
var (
	jobQueue 	chan Job
	waitGroup 	sync.WaitGroup
	workers 	[]*worker
	dispatch 	*dispatcher
)

// Worker represents the worker that executes the job
type worker struct {
	workerPool  chan chan Job
	jobChannel  chan Job
	quit    	chan bool
}

type dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool chan chan Job
	maxWorkers int
	stop    	chan bool
}

// Public methods, Start Stop

func (w *Worker) Start()  {
	jobQueue = make(chan Job, w.MaxQueue)
	dispatch = newDispatcher(w.MaxWorkers)
	dispatch.run(w.Testing)
}

func AddJob(jobs... Job)  {
	for _, job := range jobs {
		jobQueue <- job
	}
}

func GracefulShutdown()  {
	dispatch.stop <- true

	for _, worker := range workers {
		worker.stop()
	}
	waitGroup.Wait()

	workers = []*worker{}
}

// Private methods

func newDispatcher(maxWorkers int) *dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &dispatcher{workerPool: pool, maxWorkers: maxWorkers, stop: make(chan bool)}
}

func (d *dispatcher) run(t *testing.T) {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(d.workerPool)
		worker.start(t)
		workers = append(workers, &worker)
	}

	go d.dispatch()
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerPool
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		case <- d.stop:
			return
		}
	}
}

func newWorker(workerPool chan chan Job) worker {
	return worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w worker) start(t *testing.T) {
	go func(t *testing.T) {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:

				waitGroup.Add(1)

				if t == nil {
					job.perform()
				}else{
					job.performTest(t)
				}

				waitGroup.Done()

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}(t)
}

// Stop signals the worker to stop listening for work requests.
func (w worker) stop() {
	w.quit <- true
}