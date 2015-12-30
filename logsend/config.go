package logsend

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	logpkg "log"
	"os"
)

//配置结构
type Configuration struct {
	WatchDir          string
	Logger            *logpkg.Logger
	ReadWholeLog      bool
	ReadAlway         bool
	SenderName        string
	registeredSenders map[string]*SenderRegister
}

var Conf = &Configuration{
	WatchDir:          "",
	Logger:            logpkg.New(os.Stderr, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile),
	registeredSenders: make(map[string]*SenderRegister),
}

//默认配置结构
var (
	rawConfig = make(map[string]interface{}, 0)
)

//载入默认配置
func LoadRawConfig(f *flag.Flag) {
	rawConfig[f.Name] = f.Value
}

//载入自定义配置文件
func LoadConfigFromFile(fileName string) (rule *Rule, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	rawConfig, err := ioutil.ReadAll(file)
	if err != nil {
		Conf.Logger.Fatalln(err)
	}
	return LoadConfig(rawConfig)
}

//载入具体配置项
func LoadConfig(rawConfig []byte) (rule *Rule, err error) {
	config := make(map[string]interface{})
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}
	sender_map := config["senders"].(map[string]interface{})
	senders := make([]Sender, 0)
	for sender_name, register := range Conf.registeredSenders {
		if val, ok := sender_map[sender_name]; ok {
			//sender进行配置信息初始化
			register.Init(val)
			if register.initialized != true {
				continue
			}
			sender := register.get()
			senders = append(senders, sender)
		}
	}

	watch_dir := config["watchDir"].(interface{}).(string)
	regexp := config["regexp"].(interface{}).(string)

	//建立规则, 如出现异常,则立刻panic,程序终止
	rule, err = NewRule(regexp, watch_dir)
	if err != nil {
		panic(err)
	}
	rule.senders = senders
	return
}
