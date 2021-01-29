package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.0.1"
    BuildDate = "29/01/2021 at 14:37:32"
    LinesOfCode = "3020"
    NumFiles = "46"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
