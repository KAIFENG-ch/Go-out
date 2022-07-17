package Go_out

import (
	"io"
	"os"
)

type RecoveryFunc func(context *Context, err interface{})

var DefaultErrWriter io.Writer = os.Stderr

func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultWriter)
}

func CustomerRecovery(handle RecoveryFunc) HandlerFunc {
	return RecoveryWithWriter(DefaultErrWriter, handle)
}

func RecoveryWithWriter(errWriter io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	if len(recovery) > 0 {
		return nil
	}

	//TODO
	return nil
}

func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
	//var logger *log.Logger
	//if out != nil {
	//	logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	//}
	//return func(c *Context) {
	//	defer func() {
	//	}()
	//	c.Next()
	//}
	// TODO
	return nil
}