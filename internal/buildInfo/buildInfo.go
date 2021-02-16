package buildInfo

import (
    "runtime"
    "time"
)

const (
    Version = "3.3.0"
    BuildDate = "16/02/2021 at 19:53:52"
    LinesOfCode = "3233"
    NumFiles = "50"
)

var (
    GoVersion = runtime.Version()
    StartTime = time.Now()
)
