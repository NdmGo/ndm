package handles

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"ndm/internal/common"
	"ndm/internal/utils"

	"github.com/gin-gonic/gin"
)

var sysStatus struct {
	Uptime       string
	NumGoroutine int

	// General statistics.
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}
var m *runtime.MemStats

// initTime is the time when the application was initialized.
var initTime = time.Now()

func init() {
	m = new(runtime.MemStats)
}

func updateSystemStatus() {

	runtime.ReadMemStats(m)

	sysStatus.Uptime = utils.TimeSincePro(initTime)

	sysStatus.NumGoroutine = runtime.NumGoroutine()

	sysStatus.MemAllocated = utils.ToSize(int64(m.Alloc))
	sysStatus.MemTotal = utils.ToSize(int64(m.TotalAlloc))
	sysStatus.MemSys = utils.ToSize(int64(m.Sys))
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = utils.ToSize(int64(m.HeapAlloc))
	sysStatus.HeapSys = utils.ToSize(int64(m.HeapSys))
	sysStatus.HeapIdle = utils.ToSize(int64(m.HeapIdle))
	sysStatus.HeapInuse = utils.ToSize(int64(m.HeapInuse))
	sysStatus.HeapReleased = utils.ToSize(int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = utils.ToSize(int64(m.StackInuse))
	sysStatus.StackSys = utils.ToSize(int64(m.StackSys))
	sysStatus.MSpanInuse = utils.ToSize(int64(m.MSpanInuse))
	sysStatus.MSpanSys = utils.ToSize(int64(m.MSpanSys))
	sysStatus.MCacheInuse = utils.ToSize(int64(m.MCacheInuse))
	sysStatus.MCacheSys = utils.ToSize(int64(m.MCacheSys))
	sysStatus.BuckHashSys = utils.ToSize(int64(m.BuckHashSys))
	sysStatus.GCSys = utils.ToSize(int64(m.GCSys))
	sysStatus.OtherSys = utils.ToSize(int64(m.OtherSys))

	sysStatus.NextGC = utils.ToSize(int64(m.NextGC))
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
}

func SettingPage(c *gin.Context) {
	data := common.CommonVer()

	action := c.Param("action")
	data["setting_page"] = action

	updateSystemStatus()
	data["sys_state"] = sysStatus
	c.HTML(http.StatusOK, "setting.tmpl", data)
}
