package fileMonitor

import (
	"encoding/json"
	"os"
	"time"
)

type LogEntry struct {
	Action    string   `json:"action"`
	Path      string   `json:"path"`
	Timestamp string   `json:"timestamp"`
	Level     LogLevel `json:"level"`
}

// LogLevel 日志告警等级1-4从低到高
type LogLevel int32

const (
	LogLevelINFO  LogLevel = 1
	LogLevelWARN  LogLevel = 2
	LogLevelERROR LogLevel = 3
	LogLevelFATAL LogLevel = 4
)

func logJSON(logFile *os.File, action, path string, level LogLevel) error {
	// 创建一个JSON编码器
	encoder := json.NewEncoder(logFile)

	logEntry := LogEntry{
		Action:    action,
		Path:      path,
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
	}

	// 将日志写入文件
	err := encoder.Encode(logEntry)
	if err != nil {
		return err
	}

	return nil
}
