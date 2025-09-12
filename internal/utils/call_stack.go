package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func GetStackTrace(skip, depth int) []string {
	traces := make([]string, 0, depth)

	for i := skip; i < skip+depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		function := fn.Name()

		if parts := strings.Split(function, "."); len(parts) > 0 {
			traces = append(traces, fmt.Sprintf("%s:%d (%s)", file, line, parts[len(parts)-1]))
		} else {
			traces = append(traces, function)
		}
	}
	return traces
}
