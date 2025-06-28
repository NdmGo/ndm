package multitasking

import (
	// "strings"
	"sync"
)

const (
	maxConcurrent = 5 // maximum concurrent download quantity
)

var (
	instance *MultiTasking
	once     sync.Once
	wg       sync.WaitGroup
)

type DownloadTask struct {
	Key  string
	Size int64
}

type MultiTasking struct {
	limit     int64
	taskChan  chan DownloadTask
	results   chan string
	errorChan chan error
}

func Instance() *MultiTasking {
	once.Do(func() {
		instance = &MultiTasking{}
		instance.Init(10)
	})
	return instance
}

func (mt *MultiTasking) Init(limit int64) error {
	mt.limit = limit

	// create task channels and work pools
	mt.taskChan = make(chan DownloadTask, mt.limit*2)
	mt.results = make(chan string, mt.limit*2)
	mt.errorChan = make(chan error, mt.limit*2)
	return nil
}

func (mt *MultiTasking) GetTaskLimit() int64 {
	return mt.limit
}

func (mt *MultiTasking) DoneTask(fn func()) {
	wg.Add(1)
	fn()

	wg.Done()
}
