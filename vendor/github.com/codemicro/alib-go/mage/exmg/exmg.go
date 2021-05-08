package exmg

import (
	"os"
	"runtime"
)

func GetTargetOS() string {
	val, ok := os.LookupEnv("GOOS")
	if !ok {
		return runtime.GOOS
	} else {
		return val
	}
}

func GetTargetArch() string {
	val, ok := os.LookupEnv("GOARCH")
	if !ok {
		return runtime.GOARCH
	} else {
		return val
	}
}
