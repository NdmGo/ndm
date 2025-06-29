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
	limit            int64
	current_task_num int64

	taskChan  chan DownloadTask
	results   chan string
	errorChan chan error
	done      chan bool

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
	mm      map[string]*MultiTasking
)

func Instance() *MultiTasking {
	once.Do(func() {
		instance = &MultiTasking{}
		instance.Init(10)
	})
	return instance
}

func Factory() *MultiTasking {
	f := &MultiTasking{}
	f.Init(10)
	return f
}

func (mt *MultiTasking) Init(limit int64) error {
	mt.limit = limit
	mt.current_task_num = 0

	// create task channels and work pools
	mt.taskChan = make(chan DownloadTask, mt.limit*2)
	mt.results = make(chan string, mt.limit*2)
	mt.errorChan = make(chan error, mt.limit*2)
	mt.task = make(chan func(), mt.limit)

	mt.done = make(chan bool)
	go func() {
		var successCount, failCount int
		for {
			select {
			case <-mt.results:
				successCount++
			case err := <-mt.errorChan:
				log.Printf("下载失败: %v", err)
				failCount++
			case <-mt.done:
				fmt.Printf("\n下载完成! 成功: %d, 失败: %d\n", successCount, failCount)
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

func (mt *MultiTasking) do() {
	for {
		select {
		case fn := <-mt.task:
			defer wg.Done()
			mt.current_task_num -= 1
			fmt.Println("task:", len(mt.task))
			fmt.Println("do....")
			fn()

			// case err := <-mt.errorChan:
			// 	log.Printf("下载失败: %v", err)
			// 	failCount++
			// case <-mt.done:
			// 	fmt.Printf("\n下载完成! 成功: %d, 失败: %d\n", successCount, failCount)
			// 	return
		}
	}

	// for _, fn := range mt.task {
	// 	fn()
	// 	defer wg.Done()
	// 	mt.current_task_num -= 1

	// }
}

func (mt *MultiTasking) DoneTask(fn func()) {
	wg.Add(1)
	mt.current_task_num += 1
	mt.task <- fn
}

func (mt *MultiTasking) Close() {
	close(mt.taskChan)
	wg.Wait()
	close(mt.results)
	close(mt.errorChan)
	mt.done <- true
}
