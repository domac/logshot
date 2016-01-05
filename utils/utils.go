package utils

import (
	"fmt"
	"github.com/davecheney/profile"
)

//生成快照文件
func GenProfile() {
	cfg := profile.Config{
		CPUProfile:     true,
		MemProfile:     true,
		NoShutdownHook: true, // do not hook SIGINT
		ProfilePath:    ".",  //调用文件所在的目录
	}
	p := profile.Start(&cfg)
	defer p.Stop()
	fmt.Println("profile file generate success.")
}
