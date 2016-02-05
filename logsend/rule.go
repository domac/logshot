package logsend

import (
	"os"
	"regexp"
	"study2016/logshot/logger"
)

//规则引擎
type Rule struct {
	regexp   *regexp.Regexp
	watchDir string //文件监听目录
	sender   Sender //数据发送器
}

//创建规则
func NewRule(sregexp string, watchDir string) (*Rule, error) {
	rule := &Rule{}
	rule.watchDir = watchDir
	var err error

	//对watch dir 进行判断
	fi, fi_err := os.Stat(watchDir)
	if fi == nil {
		logger.GetLogger().Infoln("watch dir didn't exists !")
		return rule, fi_err
	}
	//对正则进行校验
	if rule.regexp, err = regexp.Compile(sregexp); err != nil {
		return rule, err
	}
	return rule, nil
}

//关闭Sender
func (self *Rule) CloseSender() {
	self.sender.Stop()
}

func (self *Rule) GetSender() Sender {
	return self.sender
}
