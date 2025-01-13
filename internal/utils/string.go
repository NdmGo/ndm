package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
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

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return HumanateBytes(uint64(s), 1024, sizes)
}

func ToSize(s int64) string {
	return SizeFormat(float64(s))
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
