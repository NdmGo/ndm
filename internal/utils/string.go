package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func HumanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(Logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

func SizeFormat(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	n := 0
	for size > 1024 {
		size /= 1024
		n += 1
	}

	return strconv.FormatFloat(size, 'f', 2, 32) + " " + units[n]
}

func FormatDuration(d time.Duration) string {
	// 计算天数
	days := int(d.Hours() / 24)
	if days > 0 {
		return fmt.Sprintf("%d天 %s", days, FormatDuration(d-time.Duration(days)*24*time.Hour))
	}

	// 计算小时
	hours := int(d.Hours())
	if hours > 0 {
		return fmt.Sprintf("%d小时 %s", hours, FormatDuration(d-time.Duration(hours)*time.Hour))
	}

	// 计算分钟
	minutes := int(d.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%d分钟 %s", minutes, FormatDuration(d-time.Duration(minutes)*time.Minute))
	}

	// 剩余秒数
	seconds := int(d.Seconds())
	return fmt.Sprintf("%d秒", seconds)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return HumanateBytes(uint64(s), 1024, sizes)
}

func ToSize(s int64) string {
	return SizeFormat(float64(s))
}

// TimeSincePro calculates the time interval and generate full user-friendly string.
func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func RandString(len int) string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
