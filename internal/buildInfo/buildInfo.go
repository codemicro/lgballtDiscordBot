package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.8.0"
    BuildDate = "03/01/2021 at 10:37:55"
    LinesOfCode = "2219"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
