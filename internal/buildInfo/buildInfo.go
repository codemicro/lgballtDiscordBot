package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.7.2"
    BuildDate = "02/01/2021 at 14:42:18"
    LinesOfCode = "2038"
    NumFiles = "35"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
