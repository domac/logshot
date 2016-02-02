package logsend

import (
	"os"
	"fmt"
)

//检测Agent的配置环境
func CheckAgent(configFile string) {
	fin, err := os.OpenFile(configFile, os.O_RDWR, 0644)
	defer fin.Close()
	if err != nil {
		fmt.Println("agent check fail !")
	} else {
		fmt.Println("agent check success !")
	}
}
