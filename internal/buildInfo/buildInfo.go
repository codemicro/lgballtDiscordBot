//nolint:typecheck // the embedded files can cause problems when they cannot be found, since they're not committed and created by mage in CI
package buildInfo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

//go:embed version
var Version string

//go:embed currentDate
var BuildDate string

//go:embed clocData
var jdat []byte

//go:embed changelogURL
var ChangelogURL string

var (
	GoVersion   = runtime.Version() + " " + runtime.GOOS + " " + runtime.GOARCH
	StartTime   = time.Now()
	LinesOfCode = "unknown"
	NumFiles    = "unknown"
)

// https://github.com/hhatto/gocloc/blob/ecf2a9b510f6583a05c67b8705e9bd79e8015ce1/json.go#L3-L6
type jsonLanguagesResult struct {
	Languages []clocLanguage `json:"languages"`
	Total     clocLanguage   `json:"total"`
}

// https://github.com/hhatto/gocloc/blob/ecf2a9b510f6583a05c67b8705e9bd79e8015ce1/language.go#L19-L25
type clocLanguage struct {
	Name       string `json:"name,omitempty"`
	FilesCount int32  `json:"files"`
	Code       int32  `json:"code"`
	Comments   int32  `json:"comment"`
	Blanks     int32  `json:"blank"`
}

func init() {
	var jlr jsonLanguagesResult

	_ = json.Unmarshal(jdat, &jlr)

	for _, x := range jlr.Languages {
		if strings.EqualFold(x.Name, "Go") {
			LinesOfCode = fmt.Sprint(x.Code)
			NumFiles = fmt.Sprint(x.FilesCount)
			break
		}
	}

	if BuildDate == "" {
		BuildDate = "unknown"
	} else {
		BuildDate = strings.Trim(BuildDate, "\n")
	}
}
