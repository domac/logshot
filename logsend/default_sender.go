package logsend

import (
	"fmt"
)

var (
	defaultSendCh = make(chan interface{}, 0)
)

func init() {
	RegisterNewSender("default", InitDefault, NewDefaultSender)
}

type DefaultSender struct {
	sendCh chan interface{}
}

//1.初始化配置
//2.监听消息发送通道
func InitDefault(conf interface{}) {
	go func() {
		//阻塞的方式接收defaultSendCh的消息
		for data := range defaultSendCh {
			switch msg := data.(type) {
			case *string:
				fmt.Println("standard output ===>", *msg)
			}
		}
	}()
}

//工厂类,生成本Sender
func NewDefaultSender() Sender {
	sender := &DefaultSender{
		sendCh: defaultSendCh,
	}
	return Sender(sender)
}

//注入配置
func (self *DefaultSender) SetConfig(interface{}) error {
	return nil
}

//display the name of sender
func (self *DefaultSender) Name() string {
	return "default"
}

//发送数据
func (self *DefaultSender) Send(data interface{}) {
	//不直接处理数据,先推到非缓冲的channel里面
	//fmt.Println("default send is receiving data from watch dir !")
	defaultSendCh <- data
}

func (self *DefaultSender) Stop() error {
	close(defaultSendCh)
	return nil
}
