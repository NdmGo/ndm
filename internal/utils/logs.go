package utils

import (
	"fmt"
	"os"

	"github.com/hpcloud/tail"
)

func WriteBackupLog(base_path, name, content string) error {
	l := fmt.Sprintf("%s/%s", base_path, "backup_"+name+".log")
	file, err := os.OpenFile(l, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}

func TailBackupFile(base_path, name string, n int) ([]string, error) {
	abs_path := fmt.Sprintf("%s/%s", base_path, "backup_"+name+".log")
	return TailFile(abs_path, n)
}

func TailFile(name string, n int) ([]string, error) {
	t, err := tail.TailFile(name, tail.Config{
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
