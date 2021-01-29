package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.0.0"
    BuildDate = "29/01/2021 at 09:52:20"
    LinesOfCode = "3017"
    NumFiles = "46"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
