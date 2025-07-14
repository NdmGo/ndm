package multitasking

import (
	"fmt"
	"sync"
	// "time"

	"ndm/pkg/generic_sync"

	"github.com/panjf2000/ants"
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

func Factory(name string, tag string) *MultiTasking {
	key := fmt.Sprintf("%s_%s", name, tag)

	if taskMap.Has(key) {
		task, _ := taskMap.Load(key)
		return task
	}

	f := &MultiTasking{}
	f.Init(10)
	taskMap.Store(key, f)
	return f
}

type MultiTasking struct {
	Pool    *ants.Pool
	running bool
}

func (mt *MultiTasking) Init(limit int) error {
	mt.Pool, _ = ants.NewPool(limit)
	return nil
}

func (mt *MultiTasking) IsRun() bool {
	return mt.running
}

func (mt *MultiTasking) SetForceRuningStatus() {
	mt.running = true
}

func (mt *MultiTasking) SetTaskLimit(limit uint) error {
	mt.Pool.Tune(limit)
	return nil
}

func (mt *MultiTasking) GetTaskInfo() {

}

func (mt *MultiTasking) DoneTask(callback func()) {
	mt.Pool.Submit(func() {
		callback()
	})
}

func (mt *MultiTasking) Close() {
	// defer mt.Pool.Release()
	mt.running = false
}
