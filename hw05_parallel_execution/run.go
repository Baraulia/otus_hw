package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type WorkingPool struct {
	tasks                []Task
	maxErrors            int
	tasksMU              sync.RWMutex
	numberCompletedTasks int
	errorsMU             sync.RWMutex
	numberError          int
	commonNumberTasks    int
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
		maxErrors:            m,
		tasksMU:              sync.RWMutex{},
		numberCompletedTasks: numberCompletedTask,
		errorsMU:             sync.RWMutex{},
		numberError:          numberError,
		commonNumberTasks:    len(tasks),
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
				p.tasksMU.Lock()
				p.numberCompletedTasks++
				p.tasksMU.Unlock()
				if err := task(); err != nil {
					p.errorsMU.Lock()
					p.numberError++
					p.errorsMU.Unlock()
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
		p.tasksMU.RLock()
		if p.numberCompletedTasks >= p.commonNumberTasks {
			close(p.quit)
			close(p.tasksChan)
			p.quitError = nil
			return
		}
		p.tasksMU.RUnlock()
		p.errorsMU.RLock()
		if p.numberError >= p.maxErrors {
			close(p.quit)
			close(p.tasksChan)
			p.quitError = ErrErrorsLimitExceeded
			return
		}
		p.errorsMU.RUnlock()

		p.tasksChan <- task
	}

	close(p.quit)
	close(p.tasksChan)
	p.quitError = nil
}
