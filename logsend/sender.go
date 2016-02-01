package logsend

import (
	"errors"
	"fmt"
	"study2016/logshot/logger"
)

//sender abstract
type Sender interface {
	Send(*LogLine)
	SetConfig(map[string]string) error
	Name() string
	Stop() error
}

func RegisterNewSender(name string, init func(map[string]string) error, get func() Sender) {
	sender := &SenderRegister{
		init: init,
		get:  get,
	}
	Conf.registeredSenders[name] = sender
	logger.Infoln(fmt.Sprint("register sender:", name))
	return
}

type SenderRegister struct {
	init        func(map[string]string) error
	get         func() Sender
	initialized bool
}

func (self *SenderRegister) Init(val map[string]string) error {
	err := self.init(val)
	if err != nil {
		logger.Errorln(err)
		return errors.New("sender init error")
	}
	self.initialized = true
	return nil
}
