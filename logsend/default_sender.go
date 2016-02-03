package logsend

import (
	"fmt"
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
func handleData(data *LogLine) {
	fmt.Println("[", data.Ts, "]", "standard output : ", string(data.Line))
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

func (self *DefaultSender) Receive() {
	go func() {
		for data := range self.sendCh {
			handleData(data)
		}
	}()
}
