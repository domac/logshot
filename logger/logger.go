package logger

import (
	log "github.com/cihub/seelog"
)

const (
	FATALLV = iota //0
	ERRORLV
	WARNLV
	INFOLV
	DEBUGLV
	VERBOSELV
)

//配置文件字符串形式 (主要是为了免去读取磁盘XML配置文件带来的少量IO消耗)
var configstring = `
<seelog minlevel="debug">
    <outputs>
        <filter levels="debug">
            <rollingfile formatid="sended" type="date" filename="/apps/logs/logshot/agent.log" datepattern="2006-01-02" />
        </filter>
        <filter levels="info">
            <console formatid="info"/>
            <rollingfile formatid="info" type="date" filename="/apps/logs/logshot/agent.log"  datepattern="2006-01-02" />
        </filter>
        <filter levels="warn,error,critical">
            <console formatid="error"/>
            <rollingfile formatid="error" type="date" filename="/apps/logs/logshot/agent_errors.log"  datepattern="2006-01-02" />
        </filter>
    </outputs>
    <formats>
        <format id="sended" format="%Date %Time [%Level] %Msg%n"/>
        <format id="info" format="%Date %Time [%Level] %Msg%n"/>
        <format id="error" format="%Date %Time [%Level] %Msg%n"/>
    </formats>
</seelog>
`

var sendlog *SenderLogger

type SenderLogger struct {
	LOG log.LoggerInterface
}

func GetLogger() *SenderLogger {
	return sendlog
}

func init() {
	pkglogger, err := log.LoggerFromConfigAsString(configstring)
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	//log.ReplaceLogger(pkglogger)
	sendlog = &SenderLogger{
		LOG: pkglogger,
	}
}

func (self *SenderLogger) Infoln(v ...interface{}) {
	self.LOG.Info(v)
}

func (self *SenderLogger) Infof(format string, params ...interface{}) {
	self.LOG.Infof(format, params)
}

func (self *SenderLogger) Errorln(v ...interface{}) {
	self.LOG.Error(v)
}

func (self *SenderLogger) Errorf(format string, params ...interface{}) {
	self.LOG.Errorf(format, params)
}

func (self *SenderLogger) Warnln(v ...interface{}) {
	self.LOG.Warn(v)
}

func (self *SenderLogger) Warnf(format string, params ...interface{}) {
	self.LOG.Warnf(format, params)
}

// -------  log interface ------ //

func (self *SenderLogger) Fatal(v ...interface{}) {
	self.LOG.Error(v...)
}
func (self *SenderLogger) Fatalf(format string, v ...interface{}) {
	self.LOG.Errorf(format, v...)
}
func (self *SenderLogger) Fatalln(v ...interface{}) {
	self.LOG.Error(v...)
}
func (self *SenderLogger) Panic(v ...interface{}) {
	self.LOG.Error(v...)
}
func (self *SenderLogger) Panicf(format string, v ...interface{}) {
	self.LOG.Errorf(format, v...)
}
func (self *SenderLogger) Panicln(v ...interface{}) {
	self.LOG.Error(v...)
}
func (self *SenderLogger) Print(v ...interface{}) {
	self.LOG.Info(v...)
}
func (self *SenderLogger) Printf(format string, v ...interface{}) {
	self.LOG.Infof(format, v...)
}
func (self *SenderLogger) Println(v ...interface{}) {
	self.LOG.Info(v...)
}
