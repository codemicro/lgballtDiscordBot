package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.0.0dev"
    BuildDate = "16/01/2021 at 17:39:42"
    LinesOfCode = "2491"
    NumFiles = "41"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
