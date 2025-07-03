package multitasking

import (
	"fmt"
	// "time"
	"sync"

	"ndm/pkg/generic_sync"
)

var (
	instance *MultiTasking
	once     sync.Once

	taskMap generic_sync.MapOf[string, *MultiTasking]
)

func Instance() *MultiTasking {
	once.Do(func() {
		instance = &MultiTasking{}
		instance.Init(10)
	})
	return instance
}

func Factory(name string) *MultiTasking {
	if taskMap.Has(name) {
		task, _ := taskMap.Load(name)
		return task
	}

	f := &MultiTasking{}
	f.Init(10)
	taskMap.Store(name, f)
	return f
}

type MultiTasking struct {
	limit                 int64
	current_task_num      int64
	task_already_executed int64

	backend_process bool

	running bool
	results chan string
	run_do  chan bool

	wg   sync.WaitGroup
	task chan func()
}

func (mt *MultiTasking) IsRun() bool {
	return mt.running
}

func (mt *MultiTasking) SetForceRuningStatus() {
	mt.running = true
}

func (mt *MultiTasking) do() {
	mt.wg.Add(1)
	go func() {
		defer mt.wg.Done()
		for {
			exit := false
			select {
			case fn := <-mt.task:

				fmt.Println("fn:", fn)
				fn()
				mt.current_task_num -= 1
				mt.results <- "ok"
			case <-mt.results:
				mt.task_already_executed += 1
			// case <-time.After(1 * time.Second):
			// 	fmt.Println("Working...")
			case <-mt.run_do:
				mt.running = false
				mt.backend_process = false
				exit = true
				break
			}
			if exit {
				break
			}
		}
	}()
}

func (mt *MultiTasking) reset() {
	mt.current_task_num = 0
	mt.task_already_executed = 0
}

func (mt *MultiTasking) Init(limit int64) error {

	mt.limit = limit
	mt.running = false
	mt.backend_process = true

	// create task channels and work pools
	mt.results = make(chan string, mt.limit*2)
	if mt.limit == 1 {
		mt.task = make(chan func())
	} else {
		mt.task = make(chan func(), mt.limit)
	}
	mt.run_do = make(chan bool)

	mt.reset()
	mt.do()
	return nil
}

func (mt *MultiTasking) SetTaskLimit(limit int64) error {
	mt.limit = limit
	mt.reset()
	return nil
}

func (mt *MultiTasking) GetTaskLimit() int64 {
	return mt.limit
}

func (mt *MultiTasking) GetTaskInfo() {

}

func (mt *MultiTasking) DoneTask(fn func()) {

	if !mt.running {
		mt.running = true
	}

	if !mt.backend_process {
		mt.backend_process = true
		mt.do()
	}

	mt.current_task_num += 1
	mt.task <- fn
}

func (mt *MultiTasking) Close() {
	mt.run_do <- true
	mt.wg.Wait()
	mt.reset()
}
