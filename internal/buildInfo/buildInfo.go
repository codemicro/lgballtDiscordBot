package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "2.0.1"
    BuildDate = "16/01/2021 at 12:25:28"
    LinesOfCode = "2334"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
