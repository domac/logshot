package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"study2016/logshot/logger"
	"study2016/logshot/logsend"
	"study2016/logshot/utils"
)

const (
	VERSION = "0.2.8"
)

var (
	//检测配置文件是否存在或是否定义配置文件
	check = flag.Bool("check", false, "check config.json")
	//输出Agent版本信息
	version = flag.Bool("version", false, "show version number")
	//定义发送sender
	sender = flag.String("sender", "kafka", "sender which send data to target node")
	//配置文件路径
	config = flag.String("config", "conf/config.ini", "path to config.json file")
	//读取整个日志文件
	readWholeLog = flag.Bool("readall", false, "read whole logs")
	//一直读取文件
	readAlway = flag.Bool("alway", true, "read logs once and exit")
	//是否生成性能文件
	profile = flag.Bool("profile", false, "gen profile or not")
	//定时检测
	timercheck = flag.Bool("tc", false, "timer check")
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *version {
		fmt.Printf("loghub agent version %v\n", VERSION)
		os.Exit(0)
	}

	//检测Agent
	if *check {
		logsend.CheckAgent(*config)
		os.Exit(0)
	}

	//选择sender参数
	if *sender != "" {
		logsend.Conf.SenderName = *sender
	}

	if *profile {
		utils.GenProfile()
	}

	logsend.Conf.ReadWholeLog = *readWholeLog
	logsend.Conf.ReadAlway = *readAlway

	//根据内核版本设置监听方式的配置
	if !utils.CheckKernalInotifyAbility() {
		logger.GetLogger().Infoln("watching file using polling !")
		logsend.Conf.IsPoll = true
	} else {
		logger.GetLogger().Infoln("watching file using inotify !")
	}

	//主要初始化功能服务
	logsend.InitEnv()

	fmt.Printf("Agent started, pid: %d  ppid: %d \n", os.Getpid(), os.Getppid())

	//是否开启定时检测
	if *timercheck {
		go logsend.TimerCheck()
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		logsend.WatchFiles(*config)
	} else {
		//Pipe的形式
		flag.VisitAll(logsend.LoadRawConfig)
		logsend.ProcessStdin()
	}
	os.Exit(0)
}
