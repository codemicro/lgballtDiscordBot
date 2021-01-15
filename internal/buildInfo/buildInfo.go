package buildInfo

import (
	"runtime"
	"time"
)

const (
	Version     = "2.0.0"
	BuildDate   = "15/01/2021 at 16:12:33"
	LinesOfCode = "2329"
	NumFiles    = "40"
)

var (
	GoVersion = runtime.Version()
	StartTime = time.Now()
)
