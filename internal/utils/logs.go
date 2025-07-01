package utils

import (
	// "fmt"
	// "time"

	"github.com/hpcloud/tail"
)

func WriteBackupLog() {

}

func TailFile(filename string, n int) ([]string, error) {
	t, err := tail.TailFile(filename, tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件末尾开始
		Poll:     true,
		Follow:   false,
	})
	if err != nil {
		return nil, err
	}

	var lines []string
	for line := range t.Lines {
		lines = append(lines, line.Text)
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	return lines, nil
}
