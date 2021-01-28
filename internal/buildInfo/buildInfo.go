package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.0.0"
    BuildDate = "28/01/2021 at 20:59:20"
    LinesOfCode = "3003"
    NumFiles = "46"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
