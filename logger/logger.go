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

func init() {
	logger, err := log.LoggerFromConfigAsFile("logger/seelog.xml")
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
