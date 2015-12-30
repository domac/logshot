package logsend

import (
	"os"
	"regexp"
)

//规则引擎
type Rule struct {
	regexp   *regexp.Regexp
	watchDir string   //文件监听目录
	senders  []Sender //数据发送器
	mask     string
}

//创建规则
func NewRule(sregexp string, watchDir string) (*Rule, error) {
	rule := &Rule{}
	rule.watchDir = watchDir
	var err error

	//对watch dir 进行判断
	fi, fi_err := os.Stat(watchDir)
	if fi == nil {
		Conf.Logger.Panicln("watch dir didn't exists !")
		return rule, fi_err
	}
	//对正则进行校验
	if rule.regexp, err = regexp.Compile(sregexp); err != nil {
		return rule, err
	}
	return rule, nil
}

//发送数据
func (rule *Rule) SendData(data interface{}) {
	for _, sender := range rule.senders {
		sender.Send(data)
	}
}

//关闭Sender
func (rule *Rule) CloseSender() {
	for _, sender := range rule.senders {
		sender.Stop()
	}
}
