package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.7.3"
    BuildDate = "02/01/2021 at 15:07:48"
    LinesOfCode = "2038"
    NumFiles = "35"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
