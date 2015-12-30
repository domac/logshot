package logsend

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
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

func assignFiles(allFiles []string, rule *Rule) ([]*File, error) {
	files := make([]*File, 0)
	for _, f := range allFiles {
		file, err := NewFile(f)
		if err != nil {
			return files, err
		}
		file.rule = rule
		files = append(files, file)
	}
	return files, nil
}

func continueWatch(dir *string, rule *Rule) {
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
					println("发现文件创建")
				}
			case err := <-watcher.Error:
				Conf.Logger.Println("error:", err)
			}
		}
	}()
	<-done
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

	if Conf.ReadWholeLog && Conf.ReadAlway { //全量并持续采集
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true})
	} else if Conf.ReadWholeLog { //全量但只采集一次
		file.Tail, err = tail.TailFile(fpath, tail.Config{})
	} else {
		//从当前文件最尾端开始采集
		seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Location: seekInfo})
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
	rule.SendData(line)
}

//关闭规则
func closeRule(rule *Rule) {
	rule.CloseSender()
}
