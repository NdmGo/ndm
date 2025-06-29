package multitasking

import (
	// "strings"
	"fmt"
	"log"
	"sync"
	// "time"

	"ndm/pkg/generic_sync"
)

type DownloadTask struct {
	Key  string
	Size int64
}

type MultiTasking struct {
	limit                 int64
	current_task_num      int64
	task_already_executed int64

	taskChan chan DownloadTask
	results  chan string
	done     chan bool

	task chan func()
}

const (
	maxConcurrent = 5 // maximum concurrent download quantity
)

var (
	instance *MultiTasking
	once     sync.Once
	wg       sync.WaitGroup

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

func (mt *MultiTasking) Reset() {
	mt.current_task_num = 0
	mt.task_already_executed = 0
}

func (mt *MultiTasking) Init(limit int64) error {
	mt.limit = limit
	mt.current_task_num = 0
	mt.task_already_executed = 0

	// create task channels and work pools
	mt.results = make(chan string, mt.limit*2)
	mt.task = make(chan func(), mt.limit)
	mt.done = make(chan bool)

	go func() {
		for {
			select {
			case <-mt.results:
				mt.task_already_executed += 1
				log.Printf("已经执行任务: %v", mt.task_already_executed)
			case <-mt.done:
				fmt.Println("下载完成! 成功:", mt.task_already_executed)
				return
			}
		}
	}()

	go mt.do()
	return nil
}

func (mt *MultiTasking) GetTaskLimit() int64 {
	return mt.limit
}

func (mt *MultiTasking) GetTaskInfo() {

}

func (mt *MultiTasking) do() {
	for {
		select {
		case fn := <-mt.task:
			defer wg.Done()
			mt.current_task_num -= 1
			fmt.Println("task:", len(mt.task))
			fn()
			mt.results <- "ok"
		case <-mt.done:
			fmt.Println("do over")
			return
		}
	}
}

func (mt *MultiTasking) DoneTask(fn func()) {
	wg.Add(1)
	mt.current_task_num += 1
	mt.task <- fn

}

func (mt *MultiTasking) Close() {
	wg.Wait()
	mt.done <- true
	close(mt.done)
	close(mt.results)
}
