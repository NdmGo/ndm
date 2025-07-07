package task

import (
	"context"
	"sync"
	"time"

	// "ndm/internal/conf"
	"ndm/internal/model"

	"github.com/xhofe/tache"
)

type TaskExtensionBg struct {
	tache.Base
	ctx          context.Context
	ctxInitMutex sync.Mutex
	Creator      string
	startTime    *time.Time
	endTime      *time.Time
	totalBytes   int64
}

func (t *TaskExtensionBg) SetCreator(creator string) {
	t.Creator = creator
	t.Persist()
}

func (t *TaskExtensionBg) GetCreator() string {
	return t.Creator
}

func (t *TaskExtensionBg) SetStartTime(startTime time.Time) {
	t.startTime = &startTime
}

func (t *TaskExtensionBg) GetStartTime() *time.Time {
	return t.startTime
}

func (t *TaskExtensionBg) SetEndTime(endTime time.Time) {
	t.endTime = &endTime
}

func (t *TaskExtensionBg) GetEndTime() *time.Time {
	return t.endTime
}

func (t *TaskExtensionBg) ClearEndTime() {
	t.endTime = nil
}

func (t *TaskExtensionBg) SetTotalBytes(totalBytes int64) {
	t.totalBytes = totalBytes
}

func (t *TaskExtensionBg) GetTotalBytes() int64 {
	return t.totalBytes
}

func (t *TaskExtensionBg) Ctx() context.Context {
	if t.ctx == nil {
		t.ctxInitMutex.Lock()
		if t.ctx == nil {
			t.ctx = context.WithValue(t.Base.Ctx(), "user", t.Creator)
		}
		t.ctxInitMutex.Unlock()
	}
	return t.ctx
}

func (t *TaskExtensionBg) ReinitCtx() {
	// if !conf.Tasks.AllowRetryCanceled {
	// 	return
	// }
	select {
	case <-t.Base.Ctx().Done():
		ctx, cancel := context.WithCancel(context.Background())
		t.SetCtx(ctx)
		t.SetCancelFunc(cancel)
		t.ctx = nil
	default:
	}
}

type TaskExtensionBgInfo interface {
	tache.TaskWithInfo
	GetCreator() *model.User
	GetStartTime() *time.Time
	GetEndTime() *time.Time
	GetTotalBytes() int64
}
