package logsend

import (
	"fmt"
	"os"
	"time"
	"sync"
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

var Locker sync.Mutex

func TimerCheck() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			Locker.Lock()
			fmt.Printf("the number of watching files: %d ", len(WatcherMap))
			Locker.Unlock()
		}
	}
}
