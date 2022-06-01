package Go_out

type LoggerConfig struct {
}

type LogFormatParams struct {
}

type LoggerFormatter func(params LogFormatParams) string

func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithConfig(config LoggerConfig) HandlerFunc {

	return func(context *Context) {

	}
}
