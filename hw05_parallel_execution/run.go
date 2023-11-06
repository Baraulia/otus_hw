package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type WorkingPool struct {
	tasks                []Task
	maxErrors            int32
	numberCompletedTasks int32
	numberError          int32
	commonNumberTasks    int32
	tasksChan            chan Task
	quit                 chan struct{}
	quitError            error
	wg                   *sync.WaitGroup
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var numberCompletedTask, numberError int
	pool := &WorkingPool{
		tasks:                tasks,
		maxErrors:            int32(m),
		numberCompletedTasks: int32(numberCompletedTask),
		numberError:          int32(numberError),
		commonNumberTasks:    int32(len(tasks)),
		tasksChan:            make(chan Task),
		quit:                 make(chan struct{}),
		wg:                   &sync.WaitGroup{},
	}

	if pool.maxErrors <= 0 {
		pool.maxErrors = 1
	}

	for i := 1; i <= n; i++ {
		pool.wg.Add(1)
		go pool.startWorker()
	}

	pool.wg.Add(1)
	go pool.addTasks()
	pool.wg.Wait()
	return pool.quitError
}

func (p *WorkingPool) startWorker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasksChan:
			if ok {
				atomic.AddInt32(&p.numberCompletedTasks, 1)
				if err := task(); err != nil {
					atomic.AddInt32(&p.numberError, 1)
				}
			}
		case <-p.quit:
			return
		}
	}
}

func (p *WorkingPool) addTasks() {
	defer p.wg.Done()
	for _, task := range p.tasks {
		if p.numberCompletedTasks >= p.commonNumberTasks {
			close(p.quit)
			close(p.tasksChan)
			p.quitError = nil
			return
		}

		if p.numberError >= p.maxErrors {
			close(p.quit)
			close(p.tasksChan)
			p.quitError = ErrErrorsLimitExceeded
			return
		}

		p.tasksChan <- task
	}

	close(p.quit)
	close(p.tasksChan)
	p.quitError = nil
}
