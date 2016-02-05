package logsend

import (
	"errors"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/Unknwon/com"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"strings"
	"study2016/logshot/logger"
	"time"
)

//监听的路径-文件映射
var WatcherMap map[string]*File = make(map[string]*File)

//监听文件
func WatchFiles(configFile string) {
	rule, err := LoadConfigFromFile(configFile)
	if err != nil {
		logger.GetLogger().Errorln("Can't load config", err)
	}
	files := make([]string, 0)
	files = append(files, rule.watchDir)
	assignedFiles, err := assignFiles(files, rule)
	//统计分配文件的数量
	//当每个文件完成任务后,这个数目减少1
	//主要用于判断doneCh的过程中,不至于因为doneCh的阻塞导致deadlock

	if err != nil {
		logger.GetLogger().Errorln("can't assign file per rule", err)
	}

	doneCh := make(chan string)

	//异步日志发送
	for _, file := range assignedFiles {
		file.doneCh = doneCh
		//并行处理文件采集
		go file.tail()
	}

	//如果监听对象为目录,开启异步监听
	if com.IsDir(rule.watchDir) {
		go continueWatch(&rule.watchDir, rule, doneCh)
	} else {
		go continueSingleFileWatch(&rule.watchDir, rule, doneCh)
	}

	for {
		select {
		case fpath := <-doneCh:
			if len(WatcherMap) == 0 || Conf.ReadAlway == false {
				logger.GetLogger().Infof("finished reading file %+v", fpath)
				return
			}
		}
	}
}

//单文件分配
func assignSingleFile(filepath string, rule *Rule) (*File, error) {
	is_dir := com.IsDir(filepath)
	if !is_dir {
		file, err := NewFile(filepath, rule.GetSender())
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return nil, errors.New("file not found : " + filepath)
}

//为文件分配规则
func assignFiles(allFiles []string, rule *Rule) ([]*File, error) {
	files := make([]*File, 0)
	for _, f := range allFiles {
		is_dir := com.IsDir(f)
		//watch-dir是目录的形式
		if is_dir {
			//遍历文件目录
			filepath.Walk(f, func(pth string, info os.FileInfo, err error) error {
				if err != nil {
					panic(err)
				}
				if !info.IsDir() && !strings.Contains(info.Name(), ".DS_Store") {
					file, err := NewFile(pth, rule.GetSender())
					if err != nil {
						return err
					}
					files = appendWatch(files, file)
				}
				return nil
			})
		} else {
			//分配单文件
			file, err := assignSingleFile(f, rule)
			if err != nil {
				return files, err
			}
			files = appendWatch(files, file)
		}
	}
	return files, nil
}

func appendWatch(files []*File, file *File) []*File {
	WatcherMap[file.Tail.Filename] = file
	if files != nil {
		files = append(files, file)
	}
	return files
}

//针对目录形式的监控
func continueWatch(dir *string, rule *Rule, doneCh chan string) {
	//判断dir是否是目录结构
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.GetLogger().Errorln(err.Error())
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					file, err := assignSingleFile(ev.Name, rule)
					appendWatch(nil, file)
					println("create watch file :", ev.Name, " current watch file number is:", len(WatcherMap))
					if err == nil {
						file.doneCh = doneCh
						go file.tail() //异步传输
					}
				} else if ev.IsDelete() { //文件被删除的情况,需要进行资源回收
					//获取被删除的file对象,并发完成消息到其doneCh
					logger.GetLogger().Infof("tailing file is deleted : %s ", ev.Name)
					if delete_file, ok := WatcherMap[ev.Name]; ok {
						delete(WatcherMap, ev.Name)
						if len(WatcherMap) == 0 {
							closeRule(rule)
						}
						delete_file.Tail.Stop() // stop the line tail
					} else {
						logger.GetLogger().Errorf("get delete file fail : %s", ev.Name)
					}
				}
			case err := <-watcher.Error:
				logger.GetLogger().Errorln("error:", err)
			}
		}
	}()
	//监听目录
	err = watcher.Watch(*dir)
	if err != nil {
		logger.GetLogger().Errorln(err.Error())
	}
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

//针对单文件的监控
func continueSingleFileWatch(dir *string, rule *Rule, doneCh chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.GetLogger().Errorln(err.Error())
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsDelete() { //文件被删除的情况,需要进行资源回收
					//获取被删除的file对象,并发完成消息到其doneCh
					logger.GetLogger().Infof("single tailing file is deleted : %s ", ev.Name)
					if delete_file, ok := WatcherMap[ev.Name]; ok {
						delete(WatcherMap, ev.Name)
						if len(WatcherMap) == 0 {
							closeRule(rule)
						}
						delete_file.Tail.Stop() // stop the line tail
					} else {
						logger.GetLogger().Errorf("get single delete file fail : %s", ev.Name)
					}
				}
			case err := <-watcher.Error:
				logger.GetLogger().Errorln("error:", err)
			}
		}
	}()
	//监听目录
	err = watcher.Watch(*dir)
	if err != nil {
		logger.GetLogger().Errorln(err.Error())
	}
	<-done
	/* ... do stuff ... */
	watcher.Close()
}

//监听文件结构
type File struct {
	Tail   *tail.Tail
	sender Sender
	doneCh chan string
}

//创建监听文件
func NewFile(fpath string, sender Sender) (*File, error) {
	file := &File{}
	var err error

	//是否采用低版本的poll监听方式(darwin除外)
	//Linux2.6.32以下无法使用inotity,需要把Poll打开,采用 Polling的方式
	isPoll := Conf.IsPoll
	if Conf.ReadWholeLog && Conf.ReadAlway { //全量并持续采集
		file.Tail, err = tail.TailFile(fpath, tail.Config{
			Follow: true,
			ReOpen: true,
			Poll:   isPoll,
			Logger: logger.GetLogger(), //使用自定义的日志器
		})
	} else if Conf.ReadWholeLog { //全量但只采集一次
		file.Tail, err = tail.TailFile(fpath, tail.Config{
			Poll:   isPoll,
			Logger: logger.GetLogger(), //使用自定义的日志器
		})
	} else {
		//从当前文件最尾端开始采集
		seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
		file.Tail, err = tail.TailFile(fpath, tail.Config{
			Follow:   true,
			ReOpen:   true,
			Poll:     isPoll,
			Location: seekInfo,
			Logger:   logger.GetLogger(), //使用自定义的日志器
		})
	}
	file.sender = sender
	return file, err
}

func (self *File) tail() {
	logger.GetLogger().Infoln(fmt.Sprintf("start tailing %s", self.Tail.Filename))

	//功能收尾
	defer func() {
		//通知结束通道,让主调用方结束
		logger.GetLogger().Infof("file-watching has done : %s", self.Tail.Filename)
		self.doneCh <- self.Tail.Filename
	}()
	for line := range self.Tail.Lines {
		self.sendcall(&line.Text)
	}
}

func (self *File) sendcall(line *string) {
	//日志行对象
	logline := &LogLine{
		Ts:   time.Now().UTC().UnixNano(),
		Line: []byte(*line),
	}
	self.sender.Send(logline)
}

//关闭规则
func closeRule(rule *Rule) {
	logger.GetLogger().Infoln("stop rule send channel ")
	rule.CloseSender()
}
