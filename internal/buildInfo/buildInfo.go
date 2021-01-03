package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "1.8.1"
    BuildDate = "03/01/2021 at 11:14:14"
    LinesOfCode = "2224"
    NumFiles = "40"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
