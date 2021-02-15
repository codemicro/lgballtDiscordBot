package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.2.0"
    BuildDate = "15/02/2021 at 16:05:20"
    LinesOfCode = "3266"
    NumFiles = "51"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
