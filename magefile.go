// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codemicro/alib-go/mage/exmg"
	"github.com/codemicro/alib-go/mage/exsh"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

const (
	dockerImageTag      = "lgballtdiscordbot"
	builtExecutableName = "lgballtDiscordBot"
)

var buildVersion = getVersion()

func Build() error {
	mg.Deps(InstallDeps)
	mg.Deps(PreBuild)

	fmt.Println("Building")
	_ = os.Mkdir("build", os.ModeDir)

	var fileExtension string
	if exmg.GetTargetOS() == "windows" {
		fileExtension = ".exe"
	}

	outputFilename := path.Join("build", fmt.Sprintf("%s.%s%s", builtExecutableName, buildVersion, fileExtension))


	if err := sh.Run("go", "build", "-o", outputFilename, "github.com/codemicro/lgballtDiscordBot/cmd/lgballtDiscordBot"); err != nil {
		return err
	}

	fmt.Println("Successfully built and written to", outputFilename)
	return nil
}

func InstallDeps() error {
	fmt.Println("Installing dependencies")
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	if !exsh.IsCmdAvail("gocloc") {
		fmt.Println("Installing gocloc")
		
		if err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", "-u", "github.com/hhatto/gocloc/cmd/gocloc"); err != nil {
			return err
		}

		if !exsh.IsCmdAvail("gocloc") {
			return errors.New("gocloc was installed, but cannot be found: is GOPATH/bin on PATH?")
		}

	} else {
		if mg.Verbose() {
			fmt.Println("Skipping gocloc install (found in PATH)")
		}
	}

	return nil
}

func PreBuild() error {
	fmt.Println("Running prebuild tasks")

	if !exsh.IsCmdAvail("gocloc") {
		return errors.New("gocloc must be installed on your PATH - run `mage installDeps` or see https://github.com/hhatto/gocloc")
	}

	{
		// gocloc --output-type=json . > internal/buildInfo/clocData
		gcOut, err := sh.Output("gocloc", "--output-type=json", ".")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path.Join("internal", "buildInfo", "clocData"), []byte(strings.TrimSpace(gcOut)), 0644)
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

	{
		// write internal/buildInfo/version
		err := ioutil.WriteFile(path.Join("internal", "buildInfo", "version"), []byte(buildVersion), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

type Docker mg.Namespace

func (Docker) Build() error {
	if !exsh.IsCmdAvail("docker") {
		return errors.New("docker must be installed on your PATH - see https://docs.docker.com/get-docker/")
	}

	mg.Deps(PreBuild)

	fmt.Println("Building Docker image")

	// docker build . --file Dockerfile --tag $IMAGE_NAME
	return sh.Run("docker", "build", ".", "--file", "Dockerfile", "--tag", dockerImageTag)
}

func getVersion() string {
	versionString := os.Getenv("VERSION")

	if versionString == "" {
		log.SetOutput(bytes.NewBuffer([]byte{})) // suppress mage/sh from printing the git command when run - bad solution but oh well. It works
		commitHash, err := sh.Output("git", "log", "-n1", "--format=format:'%H'")
		log.SetOutput(os.Stdout)
		if err != nil {
			return "unknown"
		}

		return strings.Trim(commitHash, "'")[:6] + "-dev"
	}

	if strings.ToLower(versionString)[0] == 'v' {
		versionString = versionString[1:]
	}

	return versionString
}
