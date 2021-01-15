package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "2.0.1"
    BuildDate = "15/01/2021 at 16:31:03"
    LinesOfCode = "2332"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
