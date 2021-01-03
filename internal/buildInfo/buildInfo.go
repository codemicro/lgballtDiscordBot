package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.8.3"
    BuildDate = "03/01/2021 at 17:36:38"
    LinesOfCode = "2225"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
