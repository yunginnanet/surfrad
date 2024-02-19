package surfrad

import (
	"fmt"
	"os"
)

var debug = false

func debugPrint(format string, args ...interface{}) {
	if debug {
		format = "DEBUG: " + format
		_, _ = os.Stderr.WriteString(fmt.Sprintf(format, args...))
	}
}
