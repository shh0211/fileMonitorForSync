package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher      *fsnotify.Watcher
	timer        *time.Timer         // 用于延时触发上传的定时器
	mu           sync.Mutex          // 用于保护共享数据的并发安全
	interval     time.Duration       // 监控文件夹的时间间隔
	isIdle       bool                // 标志文件夹是否处于稳定状态
	changedFiles map[string]struct{} // 记录本次变化的文件
}

// 创建新的 Watcher
func NewWatcher(interval time.Duration) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		watcher:      watcher,
		interval:     interval,
		isIdle:       true,
		changedFiles: make(map[string]struct{}), // 使用 map 来存储文件路径，避免重复
	}, nil
}

// 启动监控目录
func (w *Watcher) watchDir(dirPath string) error {
	// 递归监控目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// 添加目录监控
			err := w.watcher.Add(path)
			if err != nil {
				return err
			}
			log.Println("开始监控目录: ", path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 启动事件处理
	go w.handleEvents()

	return nil
}

// 处理监控事件
func (w *Watcher) handleEvents() {
	for {
		select {
		case event := <-w.watcher.Events:
			// 每当文件夹发生任何变化（创建、修改、删除、重命名），记录文件路径，并重置计时器
			log.Println("文件夹发生变化: ", event.Name)
			w.recordFileChange(event.Name)
			w.resetTimer()

		case err := <-w.watcher.Errors:
			log.Println("监控错误: ", err)
		}
	}
}

// 记录发生变化的文件路径
func (w *Watcher) recordFileChange(filePath string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.changedFiles[filePath] = struct{}{} // 将文件路径记录到 map 中
}

// 重置计时器，每次文件夹发生变化时调用
func (w *Watcher) resetTimer() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.timer != nil {
		w.timer.Stop() // 停止当前定时器
	}

	w.timer = time.AfterFunc(w.interval, func() {
		w.mu.Lock()
		w.isIdle = true
		w.mu.Unlock()

		w.triggerUpload()
	})

	w.isIdle = false
}
func (w *Watcher) triggerUpload() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isIdle {
		log.Println("文件夹未处于稳定状态，跳过上传")
		return
	}

	log.Println("10分钟内无变化，开始上传本次变化的文件...")
	for file := range w.changedFiles {
		log.Printf("上传文件: %s\n", file)
		uploadFileToLinux(file)
	}

	w.changedFiles = make(map[string]struct{})
}

func uploadFileToLinux(filepath string) {
	cmd := exec.Command("scp", filepath, fmt.Sprintf("%s@%s:%s", remoteUser, remoteHost, remoteDir))
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		println("failed to upload file: %v, output: %s", err, string(output))
	}

	println("file:" + filepath + " uploaded success")
}

var (
	watchedDir string
	remoteHost string
	remoteUser string
	remoteDir  string
)

func init() {
	flag.StringVar(&watchedDir, "watchedDir", "", "受监控目录")
	flag.StringVar(&remoteHost, "remoteHost", "", "服务器地址")
	flag.StringVar(&remoteUser, "remoteUser", "", "登录用户名")
	flag.StringVar(&remoteDir, "remoteDir", "", "服务器同步目录")
}

func main() {
	flag.Parse()
	// 检查命令行参数是否全部给出
	if watchedDir == "" || remoteHost == "" || remoteUser == "" || remoteDir == "" {
		log.Println("请提供所有必要的命令行参数")
		flag.Usage()
		os.Exit(1)
	}
	// 创建 Watcher，设置延迟时间为10分钟
	w, err := NewWatcher(10 * time.Second)
	if err != nil {
		log.Fatal("创建监视器失败:", err)
	}

	// 指定要监控的目录
	err = w.watchDir(watchedDir)
	if err != nil {
		log.Fatal("监控目录失败:", err)
	}

	select {}
}
