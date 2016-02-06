package heartbeat

import (
	"encoding/json"
	"net/http"
	"study2016/logshot/logger"
	"time"
	"io/ioutil"
	"errors"
	"fmt"
	"study2016/logshot/logsend"
)

const NotAvailableMessage = "Not available"

var CommitHash string
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

//执行心跳任务
func RunHeartBeatTask(address string) {
	http.HandleFunc("/hb", hb_hanler)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		logger.GetLogger().Errorln(err)
	}
}

func hb_hanler(rw http.ResponseWriter, r *http.Request) {
	hash := fmt.Sprintf("the number of watching files: %d \n", len(logsend.WatcherMap))
	if hash == "" {
		hash = NotAvailableMessage
	}
	uptime := time.Since(StartTime).String()
	err := json.NewEncoder(rw).Encode(HeartbeatMessage{"running", hash, uptime})
	if err != nil {
		logger.GetLogger().Errorln(err)
	}
}

func Get(address string) (HeartbeatMessage, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", address, nil)
	resp, err := client.Do(req)
	if err != nil {
		return HeartbeatMessage{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return HeartbeatMessage{}, errors.New(fmt.Sprintf("Wrong status code: %d", resp.StatusCode))
	}
	message := HeartbeatMessage{}
	err = json.Unmarshal(b, &message)
	if err != nil {
		logger.GetLogger().Errorln(err)
	}
	return message, nil
}

