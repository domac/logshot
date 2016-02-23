package logsend

import (
	"fmt"
	"runtime"
	"strconv"
	"study2016/logshot/logger"
)

func init() {
	RegisterNewSender("default", InitDefault, NewDefaultSender)
}

type DefaultSender struct {
	sendCh chan *LogLine
}

//1.初始化配置
//2.监听消息发送通道
func InitDefault(conf map[string]string, sender Sender) error {
	logger.GetLogger().Infoln("init default sender")
	sender.Receive()
	return nil
}

//工厂类,生成本Sender
func NewDefaultSender() Sender {
	sender := &DefaultSender{
		sendCh: make(chan *LogLine, 0),
	}
	return Sender(sender)
}

//处理日志数据
func handleData(w *Worker, data *LogLine) {
	fmt.Println("[", w.Name, "/", data.Ts, "]", "standard output : ", string(data.Line))
}

//注入配置
func (self *DefaultSender) SetConfig(obj interface{}) error {
	return nil
}

//display the name of sender
func (self *DefaultSender) Name() string {
	return "default"
}

func (self *DefaultSender) Stop() error {
	logger.GetLogger().Infoln("default sender stop")
	close(self.sendCh)
	return nil
}

func (self *DefaultSender) Send(ll *LogLine) {
	self.sendCh <- ll
}

type Worker struct {
	Id   int
	Name string
}

var WORKER_NUM = runtime.NumCPU()

func NewWorker(id int, name string) *Worker {
	return &Worker{
		Id:   id,
		Name: name,
	}
}

func (self *DefaultSender) Receive() {
	for i := 0; i < WORKER_NUM; i++ {
		w := NewWorker(i, "worker_"+strconv.Itoa(i))
		go consume_data(w, self.sendCh)
	}
}

func consume_data(w *Worker, jobs <-chan *LogLine) {
	for data := range jobs {
		handleData(w, data)
	}
}
