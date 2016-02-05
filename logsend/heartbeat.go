package logsend
import "time"

var StartTime time.Time

//心跳消息结构
type HeartbeatMessage struct {
	Status string `json:"status"`
	Build  string `json:"build"`
	Uptime string `json:"uptime"`
}

func init() {
	StartTime = time.Now()
}
