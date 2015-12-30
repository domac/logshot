package logsend

//sender abstract
type Sender interface {
	Send(interface{})
	SetConfig(interface{}) error
	Name() string
	Stop() error
}

func RegisterNewSender(name string, init func(interface{}), get func() Sender) {
	sender := &SenderRegister{
		init: init,
		get:  get,
	}
	Conf.registeredSenders[name] = sender
	Conf.Logger.Println("register sender:", name)
	return
}

type SenderRegister struct {
	init        func(interface{})
	get         func() Sender
	initialized bool
}

func (self *SenderRegister) Init(val interface{}) {
	self.init(val)
	self.initialized = true
}
