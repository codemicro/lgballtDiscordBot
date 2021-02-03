package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.1.0"
    BuildDate = "03/02/2021 at 20:38:25"
    LinesOfCode = "3118"
    NumFiles = "47"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
