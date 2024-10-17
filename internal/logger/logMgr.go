// internal/logger/logMgr.go
package logger

import (
	"ffmpeg_hls_go/internal/configs"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	logInstance *logrus.Logger
	once        sync.Once
)

func SetLogLevel(config *configs.Config) {
	levelMap := map[int]logrus.Level{
		0: logrus.DebugLevel,
		1: logrus.InfoLevel,
		2: logrus.WarnLevel,
		3: logrus.ErrorLevel,
		4: logrus.FatalLevel,
		5: logrus.PanicLevel,
	}

	level, ok := levelMap[config.Logging.Loglevel]
	if !ok {
		level = logrus.InfoLevel // 默认日志级别
	}

	logInstance.SetLevel(level)
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s [%s] %s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message)), nil
}

// GetLoggerInstance 返回全局唯一的 logrus 日志管理器实例
func GetLoggerInstance() *logrus.Logger {
	once.Do(func() {
		config := configs.GetConfigInstance()
		// 创建 lumberjack 实例
		logger := &lumberjack.Logger{
			Filename:   config.Logging.LogFile,
			MaxSize:    config.Logging.MaxSize,    // 每个日志文件最大 10 MB
			MaxBackups: config.Logging.MaxBackups, // 保留最近 3 个备份
			MaxAge:     config.Logging.MaxAge,     // 保留 28 天
			Compress:   config.Logging.Compress,   // 是否压缩备份
		}

		logInstance = logrus.New()

		// 创建多输出的日志写入器
		multiWriter := io.MultiWriter(os.Stdout, logger)

		logInstance.SetOutput(multiWriter)

		// 设置日志级别
		SetLogLevel(config)

		// 设置日志格式
		logInstance.SetFormatter(&CustomFormatter{})
	})
	return logInstance
}
