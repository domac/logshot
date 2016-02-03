package logsend

import (
	"errors"
	"study2016/logshot/logger"
)

//sender abstract
type Sender interface {
	Send(*LogLine)
	Receive()
	SetConfig(interface{}) error
	Name() string
	Stop() error
}

func RegisterNewSender(name string, init func(map[string]string, Sender) error, get func() Sender) {
	sender := &SenderRegister{
		init: init,
		get:  get,
	}
	Conf.registeredSenders[name] = sender
	return
}

type SenderRegister struct {
	init        func(map[string]string, Sender) error
	get         func() Sender
	initialized bool
}

//初始化配置
func (self *SenderRegister) init_receive(val map[string]string, sender Sender) error {
	err := self.init(val, sender)
	if err != nil {
		logger.GetLogger().Errorln(err)
		return errors.New("sender init_receive error")
	}
	self.initialized = true
	return nil
}

//初始化Sender
func (self *SenderRegister) InitSender(val map[string]string) (sender Sender, err error) {
	sender = self.get()
	err = self.init_receive(val, sender)
	if err != nil {
		return nil, err
	}
	if self.initialized != true {
		return nil, errors.New("sender could not be initialized !")
	}
	return sender, nil
}
