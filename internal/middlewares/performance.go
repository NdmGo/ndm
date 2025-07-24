package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// PerformanceMetrics performance metrics
type PerformanceMetrics struct {
	RequestCount    int64         `json:"request_count"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`
	MaxDuration     time.Duration `json:"max_duration"`
	MinDuration     time.Duration `json:"min_duration"`
}

// global performance metrics
var (
	globalMetrics = &PerformanceMetrics{
		MinDuration: time.Hour, // initialize to a large value
	}
)

// PerformanceMonitor performance monitoring middleware
func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// process request
		c.Next()

		// calculate processing time
		duration := time.Since(start)

		// update performance metrics
		updateMetrics(duration)

		// log slow requests
		if duration > time.Second {
			log.Warnf("Slow request: %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}

		// add response headers
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Request-ID", generateRequestID())
	}
}

// updateMetrics update performance metrics
func updateMetrics(duration time.Duration) {
	globalMetrics.RequestCount++
	globalMetrics.TotalDuration += duration
	globalMetrics.AverageDuration = globalMetrics.TotalDuration / time.Duration(globalMetrics.RequestCount)

	if duration > globalMetrics.MaxDuration {
		globalMetrics.MaxDuration = duration
	}

	if duration < globalMetrics.MinDuration {
		globalMetrics.MinDuration = duration
	}
}

// GetMetrics get performance metrics
func GetMetrics() *PerformanceMetrics {
	return globalMetrics
}

// ResetMetrics reset performance metrics
func ResetMetrics() {
	globalMetrics = &PerformanceMetrics{
		MinDuration: time.Hour,
	}
}

// generateRequestID generate request ID
func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

// RequestSizeLimit request size limit middleware
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{"error": "Request entity too large"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORS cross-origin middleware (optimized version)
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// set CORS headers
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Response-Time, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Compression response compression middleware
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if client supports compression
		acceptEncoding := c.Request.Header.Get("Accept-Encoding")

		// set compression related headers
		if acceptEncoding != "" {
			c.Header("Vary", "Accept-Encoding")
		}

		c.Next()
	}
}

// SecurityHeaders security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// security related headers
		c.Header("X-Content-Type-Options", "nosniff")
		// c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// modify CSP policy to support frontend frameworks like layui
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")

		c.Next()
	}
}
