package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.8.5"
    BuildDate = "15/01/2021 at 16:04:32"
    LinesOfCode = "2329"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
