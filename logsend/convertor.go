package logsend

import (
	"errors"
	"strconv"
)

//interface转字符串
func Ci2string(i interface{}) (o string, err error) {
	switch msg := i.(type) {
	default:
		err = errors.New("interface not a string")
	case string:
		o = msg
	case *string:
		o = *msg
	}
	return
}

//interface转浮点
func Ci2float(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a float")
	case string:
		var fl float64
		fl, err = strconv.ParseFloat(i.(string), 64)
		o = fl
	case float64:
		o = i
	}
	return
}

//interface转整数
func Ci2int(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a int")
	case string:
		var fl float64
		fl, err = strconv.ParseFloat(i.(string), 64)
		o = int64(fl)
	case float64:
		o = int64(i.(float64))
	}
	return
}
