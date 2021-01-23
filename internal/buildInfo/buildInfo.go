package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "2.0.2"
    BuildDate = "23/01/2021 at 20:42:03"
    LinesOfCode = "2401"
    NumFiles = "42"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
