package fileMonitor

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
}

// 监测文件夹（递归）
func DirMonitor(dirPath string, logFile *os.File) error {
	//新建监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	w := Watcher{
		watcher: watcher,
	}

	//递归监控目录
	err = w.watchDir(dirPath, logFile)
	if err != nil {
		return err
	}

	//循环，一直保持运行
	select {}

	return nil
}

// 递归监控目录
func (w *Watcher) watchDir(dirPath string, logFile *os.File) error {
	fi, err := os.Stat(dirPath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("not a dir path")
	}
	//通过Walk来遍历目录下的所有子目录，调用相应的函数，包括自身
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		//只需监控目录即可，目录下的文件也在监控范围内，不需要一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watcher.Add(path)
			if err != nil {
				return err
			}
			//log.Println("添加监控: ", path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Println("开始递归监控文件夹: ", dirPath)

	//另启一个goroutine来处理监控对象的事件
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				{
					if event.Has(fsnotify.Create) {
						log.Println("创建文件（夹）: ", event.Name)
						logJSON(logFile, "创建文件（夹）", event.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(event.Name)
						if err == nil && fi.IsDir() {
							w.watcher.Add(event.Name)
							//log.Println("添加监控: ", event.Name)
						}
					}
					if event.Has(fsnotify.Write) {
						log.Println("修改文件（夹）内容: ", event.Name)
						logJSON(logFile, "修改文件（夹）内容", event.Name)
					}
					if event.Has(fsnotify.Remove) {
						//当删除的是文件时，上层文件夹会发出remove事件，这里会收到1条remove事件，os.Stat会返回err；
						//当删除的是文件夹时，由于被删除的文件夹自身也在监控下，也会发送remove事件，所以会收到2条remove事件，第1条是被删除的文件夹自身发的，第2条是上层文件夹发的
						//收到第1条事件的时候，还可以使用os.Stat来获取到原文件信息，收到第2条的时候os.Stat就会返回err了

						//下面的代码是为了删除文件夹的时候只记录1次日志，使用os.Stat的返回值来做判断
						//_, err := os.Stat(event.Name)
						//if err != nil {
						//	log.Println("删除文件（夹）: ", event.Name)
						//	logJSON(logFile, "删除文件（夹）", event.Name)
						//}
						log.Println("删除文件（夹）: ", event.Name)
						logJSON(logFile, "删除文件（夹）", event.Name)
						w.watcher.Remove(event.Name)
					}
					if event.Has(fsnotify.Rename) {
						//与删除的情况不同，重命名文件夹时，只有上层文件夹会发送rename事件，所以只收到1条rename事件
						log.Println("重命名文件（夹）: ", event.Name)
						logJSON(logFile, "重命名文件（夹）", event.Name)
						//如果是目录，则移除监控
						//但这里无法使用os.Stat来判断是否是目录，重命名后，已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						w.watcher.Remove(event.Name)
					}
					if event.Has(fsnotify.Chmod) {
						log.Println("修改文件（夹）权限: ", event.Name)
						logJSON(logFile, "修改文件（夹）权限", event.Name)
					}
				}
			case err := <-w.watcher.Errors:
				{
					log.Println("Error: ", err)
					return
				}
			}
		}
	}()

	return nil
}

// 1.非递归监测文件夹 2.监测单个文件(官方不建议)
func fileMonitor(filePath string, logFile *os.File) error {
	//新建监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	w := Watcher{
		watcher: watcher,
	}

	//监控文件
	err = w.watchFile(filePath, logFile)
	if err != nil {
		return err
	}

	//循环，一直保持运行
	select {}

	return nil
}

func (w *Watcher) watchFile(filePath string, logFile *os.File) error {
	//添加要监控的文件或文件夹
	err := w.watcher.Add(filePath)
	if err != nil {
		return err
	}
	log.Println("添加监控: ", filePath)

	//另启一个goroutine来处理监控对象的事件
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				{
					//判断事件发生的类型，如下5种
					// Create 创建
					// Write 写入
					// Remove 删除
					// Rename 重命名
					// Chmod 修改权限
					if event.Has(fsnotify.Create) {
						log.Println("创建文件（夹）: ", event.Name)
						logJSON(logFile, "创建文件（夹）", event.Name)
					}
					if event.Has(fsnotify.Write) {
						log.Println("修改文件（夹）内容: ", event.Name)
						logJSON(logFile, "修改文件（夹）内容", event.Name)
					}
					if event.Has(fsnotify.Remove) {
						log.Println("删除文件（夹）: ", event.Name)
						logJSON(logFile, "删除文件（夹）", event.Name)
					}
					if event.Has(fsnotify.Rename) {
						log.Println("重命名文件（夹）: ", event.Name)
						logJSON(logFile, "重命名文件（夹）", event.Name)
					}
					if event.Has(fsnotify.Chmod) {
						log.Println("修改文件（夹）权限: ", event.Name)
						logJSON(logFile, "修改文件（夹）权限", event.Name)
					}
				}
			case err := <-w.watcher.Errors:
				{
					log.Println("Error: ", err)
					return
				}
			}
		}
	}()

	return nil
}
