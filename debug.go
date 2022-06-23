package Go_out

import (
	"fmt"
	"io"
	"os"
)

var DefaultWriter io.Writer = os.Stderr

func IsDebugging() bool {
	return false
}

func debugPrint(msg string, value ...interface{}) {
	if IsDebugging() {

	}
}

func debugPrintError(err error) {
	if err != nil {
		fmt.Fprintf(DefaultWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
}
