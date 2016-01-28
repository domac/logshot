package logsend

import (
	"errors"
	"github.com/ActiveState/tail"
	"github.com/Unknwon/com"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//监听文件
func WatchFiles(configFile string) {
	rule, err := LoadConfigFromFile(configFile)
	if err != nil {
		Conf.Logger.Fatalln("Can't load config", err)
	}
	files := make([]string, 0)
	files = append(files, rule.watchDir)
	assignedFiles, err := assignFiles(files, rule)
	//统计分配文件的数量
	//当每个文件完成任务后,这个数目减少1
	//主要用于判断doneCh的过程中,不至于因为doneCh的阻塞导致deadlock
	assignedFilesCount := len(assignedFiles)

	if err != nil {
		Conf.Logger.Fatalln("can't assign file per rule", err)
	}

	doneCh := make(chan string)
	for _, file := range assignedFiles {
		file.doneCh = doneCh
		//并行处理文件采集
		go file.tail()
	}

	//如果监听对象为目录,开启异步监听
	if com.IsDir(rule.watchDir) {
		go continueWatch(&rule.watchDir, rule, &assignedFilesCount, doneCh)
	}

	for {
		select {
		case fpath := <-doneCh:
			assignedFilesCount = assignedFilesCount - 1
			if assignedFilesCount == 0 {
				Conf.Logger.Printf("finished reading file %+v", fpath)
				return
			}
		}
	}
}

//单文件分配
func assignSingleFile(filepath string, rule *Rule) (*File, error) {
	is_dir := com.IsDir(filepath)
	if !is_dir {
		file, err := NewFile(filepath)
		if err != nil {
			return nil, err
		}
		file.rule = rule
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
					file, err := NewFile(pth)
					if err != nil {
						return err
					}
					file.rule = rule
					files = append(files, file)
				}
				return nil
			})
		} else {
			//分配单文件
			file, err := assignSingleFile(f, rule)
			if err != nil {
				return files, err
			}
			files = append(files, file)
		}
	}
	return files, nil
}

func continueWatch(dir *string, rule *Rule, totalFileCount *int, doneCh chan string) {
	//判断dir是否是目录结构
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Conf.Logger.Fatal(err)
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					*totalFileCount = *totalFileCount + 1
					println("发现文件创建:", ev.Name, "当前监控文件数:", *totalFileCount)
					file, err := assignSingleFile(ev.Name, rule)
					if err == nil {
						file.doneCh = doneCh
						go file.tail() //异步传输
					}
				}
			case err := <-watcher.Error:
				Conf.Logger.Println("error:", err)
			}
		}
	}()
	//监听目录
	err = watcher.Watch(*dir)
	if err != nil {
		Conf.Logger.Fatal(err)
	}
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

//监听文件结构
type File struct {
	Tail   *tail.Tail
	rule   *Rule
	doneCh chan string
}

//创建监听文件
func NewFile(fpath string) (*File, error) {
	file := &File{}
	var err error

	//是否采用低版本的poll监听方式
	isPoll := Conf.IsPoll

	// Config 中 设置 Poll:true 可以解决 linux 2.6.32以下的监听问题
	// 2.6.32或以上才有的是 inotity , 但2.6.32以下无法使用, 需要把 Poll打开, 采用 Polling的方式
	if Conf.ReadWholeLog && Conf.ReadAlway { //全量并持续采集
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Poll: isPoll})
	} else if Conf.ReadWholeLog { //全量但只采集一次
		file.Tail, err = tail.TailFile(fpath, tail.Config{Poll: isPoll})
	} else {
		//从当前文件最尾端开始采集
		seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Poll: isPoll, Location: seekInfo})
	}
	return file, err
}

func (self *File) tail() {
	Conf.Logger.Printf("start tailing %+v", self.Tail.Filename)
	defer func() {
		closeRule(self.rule)
		self.doneCh <- self.Tail.Filename
	}()
	for line := range self.Tail.Lines {
		checkLineRule(&line.Text, self.rule)
	}
}

//检查并进行发送
func checkLineRule(line *string, rule *Rule) {
	//日志行对象
	logline := &LogLine{
		Ts:   time.Now().UTC().UnixNano(),
		Line: []byte(*line),
	}
	rule.SendLogLine(logline)
}

//关闭规则
func closeRule(rule *Rule) {
	rule.CloseSender()
}
