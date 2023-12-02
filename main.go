package main

import (
	"FileMonitor/fileMonitor"
	"flag"
	"log"
	"os"
	"path/filepath"
)

const FileMonitorLog = "/var/log/wormhole/file_modified_log.log"
const WatchedDir = "/opt/hole" //注意监控的文件夹里不能包括FileMonitorLog本身，否则会无限循环

func main() {
	// 定义命令行参数
	logfile := flag.String("logfile", FileMonitorLog, "path to the output log file")
	watchedDir := flag.String("watchedDir", WatchedDir, "path to watched directory")
	// 解析命令行参数
	flag.Parse()
	// 输出解析后的参数值
	log.Println("logfile: ", *logfile)
	log.Println("watchedDir: ", *watchedDir)

	path, err := filepath.Abs(*watchedDir)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	// 打开日志文件以供写入，如果文件不存在则创建
	logFile, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer logFile.Close()

	//监测文件夹
	err = fileMonitor.DirMonitor(path, logFile)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
}
