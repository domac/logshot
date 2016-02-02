package logsend

import (
	"fmt"
	"study2016/logshot/logger"
)

var (
	defaultSendCh = make(chan *LogLine, 0)
)

func init() {
	RegisterNewSender("default", InitDefault, NewDefaultSender)
}

type DefaultSender struct {
	sendCh chan *LogLine
}

//1.初始化配置
//2.监听消息发送通道
func InitDefault(conf map[string]string) error {
	go func() {
		//阻塞的方式接收defaultSendCh的消息
		for data := range defaultSendCh {
			handleData(data)
		}
	}()
	return nil
}

//处理日志数据
func handleData(data *LogLine) {
	fmt.Println("[", data.Ts, "]", "standard output : ", string(data.Line))
}

//工厂类,生成本Sender
func NewDefaultSender() Sender {
	sender := &DefaultSender{
		sendCh: defaultSendCh,
	}
	return Sender(sender)
}

//注入配置
func (self *DefaultSender) SetConfig(iniConfig map[string]string) error {
	return nil
}

//display the name of sender
func (self *DefaultSender) Name() string {
	return "default"
}

func (self *DefaultSender) Send(ll *LogLine) {
	defaultSendCh <- ll
}

func (self *DefaultSender) Stop() error {
	logger.GetLogger().Infoln("kafka sender stop")
	close(defaultSendCh)
	return nil
}
