package logsend

import (
	"bufio"
	"flag"
	"os"
	"gopkg.in/inconshreveable/log15.v2"
)

//处理Pipe的标准输入命令
func ProcessStdin() error {
	var rule *Rule
	var err error
	if rawConfig["config"].(flag.Value).String() != "" {
		configFile := rawConfig["config"].(flag.Value).String()
		rule, err = LoadConfigFromFile(configFile)
		if err != nil {
			log15.Error("Can't load config", err)
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		checkLineRule(&line, rule)
	}

	return nil
}
