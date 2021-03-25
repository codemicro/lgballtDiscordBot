// +build mage

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/codemicro/alib-go/mage/exmg"
	"github.com/codemicro/alib-go/mage/exsh"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	mg.Deps(InstallDeps)
	mg.Deps(PreBuild)

	fmt.Println("Building")
	_ = os.Mkdir("build", os.ModeDir)

	var fileExtension string
	if exmg.GetTargetOS() == "windows" {
		fileExtension = ".exe"
	}

	cmd := exec.Command("go", "build", "-o", path.Join("build", "lgballtDiscordBot" + fileExtension), "github.com/codemicro/lgballtDiscordBot/cmd/lgballtDiscordBot")
	return cmd.Run()
}

func InstallDeps() error {
	fmt.Println("Installing dependencies")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

func PreBuild() error {
	fmt.Println("Running prebuild tasks")

	if !exsh.IsCmdAvail("gocloc") {
		return errors.New("gocloc must be installed on your system path - see https://github.com/hhatto/gocloc")
	}

	{
		// gocloc --output-type=json . > internal/buildInfo/jdat
		gcOut, err := sh.Output("gocloc", "--output-type=json", ".")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path.Join("internal", "buildInfo", "jdat"), []byte(strings.TrimSpace(gcOut)), 0644)
		if err != nil {
			return err
		}
	}

	{
		// date > internal/buildInfo/currentDate
		date := time.Now().Format(time.UnixDate)
		err := ioutil.WriteFile(path.Join("internal", "buildInfo", "currentDate"), []byte(date), 0644)
		if err != nil {
			return err
		}

	}

	return nil
}

type Docker mg.Namespace

func (Docker) Build() error {
	if !exsh.IsCmdAvail("docker") {
		return errors.New("docker must be installed on your system path - see https://docs.docker.com/get-docker/")
	}

	fmt.Println("Building Docker image")

	// docker build . --file Dockerfile --tag $IMAGE_NAME
	cmd := exec.Command("docker", "build", ".", "--file", "Dockerfile", "--tag", "lgballtdiscordbot")
	return cmd.Run()
}
