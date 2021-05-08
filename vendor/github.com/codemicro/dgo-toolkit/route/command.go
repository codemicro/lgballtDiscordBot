package route

import "regexp"

// Command represents a command
type Command struct {
	Name             string
	Help             string
	CommandText      []string
	Arguments        []Argument
	Restrictions     []CommandRestriction
	Run              MessageRunFunc
	Invisible        bool // prevent command from being shown in *kit.GetCommandInfo
	AllowOverloading bool
	Category         uint

	detectRegexp *regexp.Regexp
}
