package multitasking

import (
	"fmt"
	"sync"

	"ndm/pkg/generic_sync"
)

const (
	maxConcurrent = 5 // maximum concurrent download quantity
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

	running bool
	results chan string

	wg   sync.WaitGroup
	task chan func()
}

func (mt *MultiTasking) IsRun() bool {
	if mt.running {
		return true
	}
	return false
}

func (mt *MultiTasking) do() {
	go func() {
		for {
			select {
			case fn := <-mt.task:
				mt.wg.Done()
				mt.current_task_num -= 1
				fn()
				mt.results <- "ok"
			case <-mt.results:
				mt.task_already_executed += 1
			}
		}
	}()
}

func (mt *MultiTasking) Reset() {
	mt.current_task_num = 0
	mt.task_already_executed = 0

	// create task channels and work pools
	mt.results = make(chan string, mt.limit*2)
	mt.task = make(chan func(), mt.limit)
}

func (mt *MultiTasking) Init(limit int64) error {
	mt.limit = limit
	mt.running = false
	mt.Reset()
	mt.do()

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
	mt.wg.Add(1)
	mt.current_task_num += 1
	mt.task <- fn
}

func (mt *MultiTasking) Close() {
	mt.wg.Wait()
	mt.running = false
	close(mt.task)
	mt.Reset()
}
