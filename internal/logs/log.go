package logs

import (
	"fmt"
	"os"
	"strings"

	go_logger "github.com/phachon/go-logger"

	"ndm/internal/conf"
)

var (
	logFileName = "ndm.log"
	logger      *go_logger.Logger
)

func InitLog() error {
	logger = go_logger.NewLogger()

	jsonFormat := false
	if strings.EqualFold(conf.Log.Format, "json") {
		jsonFormat = true
	}

	logPath := conf.Log.RootPath

	// creare logs dir
	err := os.MkdirAll(logPath, 0o777)
	if err != nil {
		return err
	}

	fileConfig := &go_logger.FileConfig{
		Filename: fmt.Sprintf("%s/%s", logPath, logFileName),
		LevelFileName: map[int]string{
			logger.LoggerLevel("error"): fmt.Sprintf("%s/%s", logPath, "ndm.log"),
			logger.LoggerLevel("debug"): fmt.Sprintf("%s/%s", logPath, "ndm.log"),
		},
		MaxSize:    1024 * 1024,
		MaxLine:    100000,
		DateSlice:  "d",
		JsonFormat: jsonFormat,
		Format:     "",
	}
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)

	return nil
}

func GetLogger() *go_logger.Logger {
	return logger
}

func Debug(args string) {
	logger.Debug(args)
}

func Info(args string) {
	logger.Info(args)
}

func Warn(args string) {
	logger.Warning(args)
}

func Error(args string) {
	logger.Error(args)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
