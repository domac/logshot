package main

import (
	"flag"
	"fmt"
	logpkg "log"
	"os"
	"study2016/logshot/logsend"
)

const (
	VERSION = "0.2.5"
)

var (
	//检测配置文件是否存在或是否定义配置文件
	check = flag.Bool("check", false, "check config.json")

	//输出Agent版本信息
	version = flag.Bool("version", false, "show version number")

	//应用自身日志输出文件
	logFile = flag.String("log", "/apps/logs/logshot.log", "log file")

	//定义发送sender
	sender = flag.String("sender", "default", "sender which send data to target node")

	//配置文件路径
	config = flag.String("config", "config.json", "path to config.json file")

	//读取整个日志文件
	readWholeLog = flag.Bool("readall", false, "read whole logs")

	readAlway = flag.Bool("alway", true, "read logs once and exit")
)

func main() {

	flag.Parse()

	if *version {
		fmt.Printf("loghub agent version %v\n", VERSION)
		os.Exit(0)
	}

	//自身日志输出
	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Failed to open log file: %+v\n", err)
		}
		defer file.Close()
		logsend.Conf.Logger = logpkg.New(file, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile)
	}

	//配置检查命令
	if *check {
		//载入配置文件
		_, err := logsend.LoadConfigFromFile(*config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ok")
		os.Exit(0)
	}

	//选择sender参数
	if *sender != "" {
		logsend.Conf.SenderName = *sender
	}

	logsend.Conf.ReadWholeLog = *readWholeLog
	logsend.Conf.ReadAlway = *readAlway

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		logsend.WatchFiles(*config)
	}

}
