package logsend

import (
	"bytes"
	"fmt"
	log "gopkg.in/inconshreveable/log15.v2"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat     = "2006-01-02 15:04:05"
	termTimeFormat = "01-02|15:04:05"
	floatFormat    = 'f'
	termMsgJust    = 400
)

func AgentFormat() log.Format {
	return log.FormatFunc(func(r *log.Record) []byte {
		b := &bytes.Buffer{}
		lvl := strings.ToUpper(r.Lvl.String())
		fmt.Fprintf(b, "[%s] [%s] %s ", lvl, r.Time.Format(timeFormat), r.Msg)
		// try to justify the log output for short messages
		if len(r.Ctx) > 0 && len(r.Msg) < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-len(r.Msg)))
		}
		// print the keys logfmt style
		logfmt(b, r.Ctx)
		return b.Bytes()
	})
}

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

func logfmt(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := ctx[i].(string)
		v := formatLogfmtValue(ctx[i+1])
		if !ok {
			k, v = "LOG15_ERROR", formatLogfmtValue(k)
		}

		// XXX: we should probably check that all of your key bytes aren't invalid
		fmt.Fprintf(buf, "%s=%s", k, v)
	}

	buf.WriteByte('\n')
}

func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return escapeString(v)
	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

func escapeString(s string) string {
	needQuotes := false
	e := bytes.Buffer{}
	e.WriteByte('"')
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needQuotes = true
		}

		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteByte('\\')
			e.WriteByte('n')
		case '\r':
			e.WriteByte('\\')
			e.WriteByte('r')
		case '\t':
			e.WriteByte('\\')
			e.WriteByte('t')
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	start, stop := 0, e.Len()
	if !needQuotes {
		start, stop = 1, stop-1
	}
	return string(e.Bytes()[start:stop])
}
