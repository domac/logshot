package logsend
import "github.com/juju/errors"

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
	Conf.Logger.Println("register sender:", name)
	return
}

type SenderRegister struct {
	init        func(map[string]string) error
	get         func() Sender
	initialized bool
}

func (self *SenderRegister) Init(val map[string]string) error {
	err := self.init(val)
	if err!= nil {
		return errors.New("sender init error")
	}
	self.initialized = true
	return nil
}
