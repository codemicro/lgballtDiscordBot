package exsh

import (
	"os/exec"
)

func IsCmdAvail(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
