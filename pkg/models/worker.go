package models

import (
	"fmt"
	"sync"
	"time"
)

// WokerTask encapsulates a work item that should go in a work
type WokerTask struct {
	Err  error
	Data interface{}
	f    func(interface{}) error
}

// NewWokerTask initializes a new task based on a given work function.
func NewWokerTask(f func(interface{}) error, data interface{}) *WokerTask {
	return &WokerTask{f: f, Data: data}
}

func process(workerID int, task *WokerTask) {
	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.f(task.Data)
}

// Worker handles all the work
type Worker struct {
	ID       int
	taskChan chan *WokerTask
	quit     chan bool
}

// NewWorker returns new instance of worker
func NewWorker(channel chan *WokerTask, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
	}
}

// Start starts the worker
func (wr *Worker) Start(wg *sync.WaitGroup) {
	fmt.Printf("Starting worker %d\n", wr.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range wr.taskChan {
			process(wr.ID, task)
		}
	}()
}

// StartBackground starts the worker in background waiting
func (wr *Worker) StartBackground() {
	fmt.Printf("Starting worker %d\n", wr.ID)

	for {
		select {
		case task := <-wr.taskChan:
			process(wr.ID, task)
		case <-wr.quit:
			return
		}
	}
}

// Stop quits the worker
func (wr *Worker) Stop() {
	fmt.Printf("Closing worker %d\n", wr.ID)
	go func() {
		wr.quit <- true
	}()
}

// Pool is the worker pool
type Pool struct {
	WokerTasks []*WokerTask
	Workers    []*Worker

	concurrency   int
	collector     chan *WokerTask
	runBackground chan bool
	wg            sync.WaitGroup
}

// NewPool initializes a new pool with the given tasks and
// at the given concurrency.
func NewPool(tasks []*WokerTask, concurrency int) *Pool {
	return &Pool{
		WokerTasks:  tasks,
		concurrency: concurrency,
		collector:   make(chan *WokerTask, 1000),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		worker.Start(&p.wg)
	}

	for i := range p.WokerTasks {
		p.collector <- p.WokerTasks[i]
	}
	close(p.collector)

	p.wg.Wait()
}

// AddWokerTask adds a task to the pool
func (p *Pool) AddWokerTask(task *WokerTask) {
	p.collector <- task
}

// RunBackground runs the pool in background
func (p *Pool) RunBackground() {
	go func() {
		for {
			fmt.Print("⌛ Waiting for tasks to come in ...\n")
			time.Sleep(10 * time.Second)
		}
	}()

	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.WokerTasks {
		p.collector <- p.WokerTasks[i]
	}

	p.runBackground = make(chan bool)
	<-p.runBackground
}

// Stop stops background workers
func (p *Pool) Stop() {
	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}
