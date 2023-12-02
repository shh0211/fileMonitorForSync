package fileMonitor

import (
	"encoding/json"
	"os"
	"time"
)

type LogEntry struct {
	Action    string `json:"action"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

func logJSON(logFile *os.File, action, path string) error {
	// 创建一个JSON编码器
	encoder := json.NewEncoder(logFile)

	logEntry := LogEntry{
		Action:    action,
		Path:      path,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 将日志写入文件
	err := encoder.Encode(logEntry)
	if err != nil {
		return err
	}

	return nil
}
