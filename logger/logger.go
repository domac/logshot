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

//配置文件字符串形式
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

func init() {


	logger, err := log.LoggerFromConfigAsString(configstring)
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
}

func Infoln(v ...interface{}) {
	log.Info(v)
}

func Infof(format string, params ...interface{}) {
	log.Infof(format, params)
}

func Errorln(v ...interface{}) {
	log.Error(v)
}

func Errorf(format string, params ...interface{}) {
	log.Errorf(format, params)
}

func Warnln(v ...interface{}) {
	log.Warn(v)
}

func Warnf(format string, params ...interface{}) {
	log.Warnf(format, params)
}
