package exsh

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func IsCmdAvail(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func EnsureGoBin(binaryName, importPath string) error {
	if !IsCmdAvail(binaryName) {

		if mg.Verbose() {
			fmt.Printf("Installing %s\n", binaryName)
		}

		if err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", "-u", importPath); err != nil {
			return err
		}

		if !IsCmdAvail(binaryName) {
			return fmt.Errorf("%s was installed, but cannot be found: is GOPATH/bin on PATH?", binaryName)
		}

	} else {
		if mg.Verbose() {
			fmt.Printf("Skipping %s install (found in PATH)\n", binaryName)
		}
	}

	return nil
}
