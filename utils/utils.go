package utils

import (
	"fmt"
	"github.com/davecheney/profile"
	"os/exec"
	"strconv"
	"strings"
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

//检查操作系统内核是否具备inotify的能力
func CheckKernalInotifyAbility() bool {
	inotify_version := 31
	kenel_perfix := "2.6."
	get_version_cmd := "uname -a | awk '{split($3,a,\"-\");print a[1]}'"
	cmd := exec.Command("/bin/sh", "-c", get_version_cmd)
	stdout, err := cmd.Output()
	if err != nil {
		return false
	}
	out := string(stdout)
	if strings.Contains(out, kenel_perfix) {
		substrs := strings.Split(out, kenel_perfix)
		if len(substrs) > 0 {
			version_str := substrs[1]
			version, err := strconv.Atoi(version_str)
			if err != nil {
				return false
			}
			if version <= inotify_version {
				return false
			}
		}
	}
	return true
}
