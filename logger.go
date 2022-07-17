package Go_out

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LoggerConfig struct {
	Formatter LoggerFormatter
	output io.Writer
	skipPath []string
}

type LogFormatParams struct {
	Request http.Request
	TimeStamp time.Time

	path string
	Latency time.Duration

	isTerm bool

	method string

	ClientIP string

	Keys map[string]interface{}

	StatusCode int
}

type LoggerFormatter func(params LogFormatParams) string

func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

var defaultLogFormatter = func(param LogFormatParams) string { return "" }

func LoggerWithConfig(config LoggerConfig) HandlerFunc {
	formatter := config.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := config.output
	if out == nil {
		out = DefaultWriter
	}

	notlogged := config.skipPath

	isTerm := true

	if _, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *Context) {
		start := time.Now()
		path := c.request.URL.Path
		raw := c.request.URL.RawQuery

		c.Next()

		if _, ok := skip[path]; !ok {
			param := LogFormatParams{
				Request: *c.request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}

			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.method = c.request.Method
			param.StatusCode = c.Writer.Status()

			if raw != "" {
				path = path + "?" + raw
			}

			param.path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
