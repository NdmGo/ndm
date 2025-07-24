package conf

import (
	"runtime"
	"time"
)

// PerformanceConfig performance optimization configuration
type PerformanceConfig struct {
	// database connection pool configuration
	DB struct {
		MaxIdleConns    int           // maximum idle connections
		MaxOpenConns    int           // maximum open connections
		ConnMaxLifetime time.Duration // connection maximum lifetime
		PrepareStmt     bool          // enable prepared statements
	}

	// HTTP server configuration
	HTTP struct {
		MaxConcurrent int           // maximum concurrent connections
		ReadTimeout   time.Duration // read timeout
		WriteTimeout  time.Duration // write timeout
		IdleTimeout   time.Duration // idle timeout
	}

	// cache configuration
	Cache struct {
		UserCacheExpiry time.Duration // user cache expiry time
		MaxCacheSize    int           // maximum cache size
	}

	// file processing configuration
	File struct {
		BufferSize      int // file read/write buffer size
		LargeBufferSize int // large file buffer size
		MaxMemoryBuffer int // maximum memory buffer
	}

	// concurrency control
	Concurrency struct {
		WorkerPoolSize int // worker pool size
		QueueSize      int // queue size
	}
}

// Performance global performance configuration instance
var Performance = &PerformanceConfig{}

// InitPerformanceConfig initialize performance configuration
func InitPerformanceConfig() {
	// dynamically adjust configuration based on CPU core count
	cpuCount := runtime.NumCPU()

	// database configuration
	Performance.DB.MaxIdleConns = 10
	Performance.DB.MaxOpenConns = cpuCount * 10
	Performance.DB.ConnMaxLifetime = time.Hour
	Performance.DB.PrepareStmt = true

	// HTTP configuration
	Performance.HTTP.MaxConcurrent = 1000
	Performance.HTTP.ReadTimeout = 30 * time.Second
	Performance.HTTP.WriteTimeout = 30 * time.Second
	Performance.HTTP.IdleTimeout = 60 * time.Second

	// cache configuration
	Performance.Cache.UserCacheExpiry = 5 * time.Minute
	Performance.Cache.MaxCacheSize = 1000

	// file processing configuration
	Performance.File.BufferSize = 64 * 1024             // 64KB
	Performance.File.LargeBufferSize = 1024 * 1024      // 1MB
	Performance.File.MaxMemoryBuffer = 32 * 1024 * 1024 // 32MB

	// concurrency control
	Performance.Concurrency.WorkerPoolSize = cpuCount * 2
	Performance.Concurrency.QueueSize = 1000
}

// GetOptimalBufferSize get optimal buffer size based on file size
func GetOptimalBufferSize(fileSize int64) int {
	if fileSize > 10*1024*1024 { // use large buffer for files larger than 10MB
		return Performance.File.LargeBufferSize
	}
	return Performance.File.BufferSize
}

// GetOptimalWorkerCount get optimal worker goroutine count
func GetOptimalWorkerCount() int {
	return Performance.Concurrency.WorkerPoolSize
}
