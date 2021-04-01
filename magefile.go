// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

const (
	dockerImageTag      = "lgballtdiscordbot"
	builtExecutableName = "lgballtDiscordBot"
)

var buildVersion = getLatestCommitHash(true)

func SetBuildVersion(ver string) {
	buildVersion = strings.TrimPrefix(ver, "v")
}

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

	mg.Deps(EnsureGocloc)

	return nil
}

func EnsureGocloc() error {
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

	mg.Deps(EnsureGocloc)
	mg.Deps(PreBuild)

	fmt.Println("Building Docker image as version", buildVersion)

	// docker build . --file Dockerfile --tag $IMAGE_NAME
	return sh.Run("docker", "build", ".", "--file", "Dockerfile", "--tag", dockerImageTag)
}

func (Docker) Login(registry, user string) error {
	// Requires token to be in REGISTRY_AUTH_TOKEN

	token, ok := os.LookupEnv("REGISTRY_AUTH_TOKEN")
	if !ok {
		return errors.New("REGISTRY_AUTH_TOKEN not set")
	}

	var output io.Writer
	if mg.Verbose() {
		output = os.Stdout
	} else {
		output = bytes.NewBuffer([]byte{})
	}

	fmt.Printf("Logging into Docker registry %s", registry)

	// echo "$REGISTRY_AUTH_TOKEN" | docker login ghcr.io -u codemicro --password-stdin
	cmd := exec.Command("docker", "login", registry, "-u", user, "--password-stdin")
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewBufferString(token)
	return cmd.Run()
}

func (Docker) Publish(imageId string) error {
	// note: imageId should be something like "blah/blah/blah:latest"

	imageId = fmt.Sprintf(strings.ToLower(imageId), dockerImageTag)

	fmt.Println("Publishing Docker image as", imageId)

	// docker tag $IMAGE_NAME $IMAGE_ID
	if err := sh.Run("docker", "tag", dockerImageTag, imageId); err != nil {
		return err
	}

	// docker push $IMAGE_ID
	return sh.Run("docker", "push", imageId)
}

func getLatestCommitHash(trim bool) string {

	// suppress mage/sh from printing the git command when run - bad solution but oh well. It works
	// https://github.com/magefile/mage/issues/291
	log.SetOutput(bytes.NewBuffer([]byte{})) 
	commitHash, err := sh.Output("git", "log", "-n1", "--format=format:'%H'")
	log.SetOutput(os.Stdout)
	if err != nil {
		return "unknown"
	}

	cutStr := strings.Trim(commitHash, "'")
	if trim {
		cutStr = cutStr[:6]
	}

	return cutStr + "-dev"
}
