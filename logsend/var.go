package logsend

import (
	"os"
	"study2016/logshot/logger"
	"study2016/logshot/utils"
)

var Root string       //根目录路径
var LocalIps []string //本地IP

func InitEnv() {
	initRootDir()
	InitLocalIps()
}

//初始化Agent根目录路径
func initRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		logger.GetLogger().Errorln(err.Error())
	}
}

//初始化本地IP
func InitLocalIps() {
	var err error
	LocalIps, err = utils.IntranetIP()
	if err != nil {
		logger.GetLogger().Fatalln("get intranet ip fail:", err)
	}
}
