package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/hpcloud/tail"
)

func TruncateLog(base_path, prefix, name string) error {
	l := fmt.Sprintf("%s/%s_%s", base_path, prefix, name+".log")
	err := os.Truncate(l, 0)
	if err != nil {
		return err
	}
	return nil
}

func WriteLog(base_path, prefix, name, content string) error {
	l := fmt.Sprintf("%s/%s_%s", base_path, prefix, name+".log")
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

func TailFile(base_path, prefix, name string, n int) ([]string, error) {
	abs_path := fmt.Sprintf("%s/%s_%s", base_path, prefix, name+".log")
	if !IsExist(abs_path) {
		return []string{}, nil
	}
	return getLastNLines(abs_path, n)
}

func getLastNLines(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(lines) >= n {
			lines = append(lines[1:], scanner.Text())
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func readLastLines(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, _ := file.Stat()
	filesize := stat.Size()
	var cursor int64 = 0
	lines := make([]string, 0, n)
	line := ""

	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		_, err := file.Read(char)
		if err != nil {
			break
		}

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // 换行或回车
			if len(line) > 0 {
				lines = append(lines, line)
				line = ""
				if len(lines) == n {
					break
				}
			}
		} else {
			line = string(char) + line
		}

		if cursor == -filesize { // 到达文件开头
			if len(line) > 0 {
				lines = append(lines, line)
			}
			break
		}
	}

	// 反转顺序
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines, nil
}

func GetLastNLinesSeek(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := stat.Size()

	var lines []string
	buf := make([]byte, 1024) // 缓冲区大小
	// lineBreak := []byte{'\n'}
	pos := size

	for len(lines) < n && pos > 0 {
		// 计算读取位置和大小
		readSize := int64(len(buf))
		if pos < readSize {
			readSize = pos
		}
		pos -= readSize

		// 读取块
		_, err := file.Seek(pos, io.SeekStart)
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(file, buf[:readSize])
		if err != nil {
			return nil, err
		}

		// 在块中查找换行符
		chunk := buf[:readSize]
		for i := len(chunk) - 1; i >= 0; i-- {
			if chunk[i] == '\n' {
				line := string(chunk[i+1:])
				lines = append(lines, line)
				if len(lines) == n {
					break
				}
			}
		}
	}

	// 反转顺序
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines, nil
}

func TailFileBak1(name string, n int) ([]string, error) {
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
