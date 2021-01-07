package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.8.4"
    BuildDate = "07/01/2021 at 12:01:41"
    LinesOfCode = "2237"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
