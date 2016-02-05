package logsend

import (
	"bufio"
	"flag"
	"os"
	"study2016/logshot/logger"
	"time"
)

//处理Pipe的标准输入命令
func ProcessStdin() error {
	var rule *Rule
	var err error
	if rawConfig["config"].(flag.Value).String() != "" {
		configFile := rawConfig["config"].(flag.Value).String()
		rule, err = LoadConfigFromFile(configFile)
		if err != nil {
			logger.GetLogger().Errorln("Can't load config", err)
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		logline := &LogLine{
			Ts:   time.Now().UTC().UnixNano(),
			Line: []byte(line),
		}
		rule.GetSender().Send(logline)
	}

	return nil
}
