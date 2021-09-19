// +build mage

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

	_ = os.Mkdir("build", os.ModeDir)

	var fileExtension string
	if exmg.GetTargetOS() == "windows" {
		fileExtension = ".exe"
	}

	outputFilename := path.Join("build", fmt.Sprintf("%s.%s%s", builtExecutableName, buildVersion, fileExtension))

	if err := sh.Run("go", "build", "-o", outputFilename, "github.com/codemicro/lgballtDiscordBot/cmd/lgballtDiscordBot"); err != nil {
		return err
	}

	fmt.Println(outputFilename)
	return nil
}

func InstallDeps() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	mg.Deps(VendorDeps)
	mg.Deps(EnsureGocloc)
	mg.Deps(EnsureQTC)

	return nil
}

func VendorDeps() error {
	if err := sh.Run("go", "mod", "vendor"); err != nil {
		return err
	}
	// running `go mod vendor` undoes any applied patches
	mg.Deps(ApplyPatches)

	return nil
}

func EnsureGocloc() error {
	if !exsh.IsCmdAvail("gocloc") {
		if err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", "-u", "github.com/hhatto/gocloc/cmd/gocloc"); err != nil {
			return err
		}

		if !exsh.IsCmdAvail("gocloc") {
			return errors.New("gocloc was installed, but cannot be found: is GOPATH/bin on PATH?")
		}

	} else {
		if mg.Verbose() {
			fmt.Fprintln(os.Stderr, "Skipping gocloc install (found in PATH)")
		}
	}

	return nil
}

func EnsureQTC() error {
	if !exsh.IsCmdAvail("qtc") {
		if err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", "-u", "github.com/valyala/quicktemplate/qtc"); err != nil {
			return err
		}

		if !exsh.IsCmdAvail("qtc") {
			return errors.New("qtc was installed, but cannot be found: is GOPATH/bin on PATH?")
		}

	} else {
		if mg.Verbose() {
			fmt.Fprintln(os.Stderr, "Skipping qtc install (found in PATH)")
		}
	}

	return nil
}

func ApplyPatches() error {
	// read patches index
	patches := make(map[string]string)
	{
		dat, err := ioutil.ReadFile("patches/patches.json")
		if err != nil {
			return err
		}
		err = json.Unmarshal(dat, &patches)
		if err != nil {
			return err
		}
	}

	// apply patches
	for patch, directory := range patches {

		patchPath := strings.Join([]string{"patches", patch}, string(os.PathSeparator))

		err := sh.Run("git", "apply", "--directory="+directory, "--ignore-space-change", "--ignore-whitespace", patchPath)
		if err != nil {
			fmt.Printf("WARNING: Failed to apply Git patch %s (%s)\n", patch, err.Error())
		}
	}

	return nil
}

func PreBuild() error {
	if !exsh.IsCmdAvail("gocloc") {
		return errors.New("gocloc must be installed on your PATH - run `mage installDeps` or see https://github.com/hhatto/gocloc")
	}

	if !exsh.IsCmdAvail("qtc") {
		return errors.New("qtc must be installed on your PATH - run `mage installDeps` or see https://github.com/valyala/quicktemplate")
	}

	{
		// gocloc --output-type=json --not-match-d=vendor . > internal/buildInfo/clocData
		gcOut, err := sh.Output("gocloc", "--output-type=json", "--not-match-d=vendor", ".")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path.Join("internal", "buildInfo", "clocData"), []byte(strings.TrimSpace(gcOut)), 0644)
		if err != nil {
			return err
		}
	}

	{
		// qtc -dir=internal/adminSite/templates -skipLineComments
		err := sh.Run("qtc", "-dir=internal/adminSite/templates", "-skipLineComments")
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

	{
		// upload changelog
		changelogBytes, err := ioutil.ReadFile(path.Join(".github", "CHANGELOG.md"))
		if err != nil {
			return err
		}
		changelogString := string(changelogBytes)
		changelogString = strings.ReplaceAll(changelogString, "codemicro", "username-censored")

		// based on this: https://github.com/radude/rentry/blob/master/rentry

		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}

		client := &http.Client{
			Jar: jar,
		}

		csrfResponse, err := client.Get("https://rentry.co")
		if err != nil {
			return err
		}
		if csrfResponse.StatusCode != 200 {
			return fmt.Errorf("unknown response code of %d from Rentry", csrfResponse.StatusCode)
		}

		u, _ := url.Parse("https://rentry.co")
		var v string
		for _, cook := range jar.Cookies(u) {
			if cook.Name == "csrftoken" {
				v = cook.Value
				break
			}
		}

		form := url.Values{}
		form.Add("text", changelogString)
		form.Add("csrfmiddlewaretoken", v)

		req, err := http.NewRequest("POST", "https://rentry.co/api/new", strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Referer", "https://rentry.co")

		resp, err := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("unknown response code of %d from Rentry\n%s", resp.StatusCode, string(body))
		}

		output := struct {
			Status   string `json:"status"`
			Content  string `json:"content"`
			Url      string `json:"url"`
			EditCode string `json:"edit_code"`
		}{}

		err = json.Unmarshal(body, &output)
		if err != nil {
			return err
		}

		// write internal/buildInfo/changelogURL
		err = ioutil.WriteFile(path.Join("internal", "buildInfo", "changelogURL"), []byte(output.Url), 0644)
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

	if mg.Verbose() {
		fmt.Fprintln(os.Stderr, "Building Docker image as version", buildVersion)
	}

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

	if mg.Verbose() {
		fmt.Fprintln(os.Stderr, "Publishing Docker image as", imageId)
	}

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
