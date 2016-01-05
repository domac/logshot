package logsend

import (
	"encoding/json"
)

type LogLine struct {
	Ts   int64
	Line []byte
}

func MarshaLogLines(loglines []*LogLine) []byte {
	b, err := json.Marshal(loglines)
	if err != nil {
		panic(err)
	}
	return b
}
