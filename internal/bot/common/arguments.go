package common

import "regexp"

var spaceSplitRegex = regexp.MustCompile(` +`)

func TakeFirstPart(in string) (string, string) {
	xspl := spaceSplitRegex.Split(in, 2)
	var v string
	if len(xspl) > 1 {
		v = xspl[1]
	}
	return xspl[0], v
}

var idFromPingRegex = regexp.MustCompile(`<@!?(.+)>`)

var PingOrUserIdType = pingOrUserIdType{}

type pingOrUserIdType struct{}

func (pingOrUserIdType) Parse(content *string) (interface{}, error) {

	a, b := TakeFirstPart(*content)

	var userId string

	if idFromPingRegex.MatchString(a) {
		matches := idFromPingRegex.FindStringSubmatch(a)
		userId = matches[1]
	} else {
		userId = a
	}

	*content = b
	return userId, nil

}
func (pingOrUserIdType) Help(_ string) string { return "A user ID or ping" }
func (pingOrUserIdType) Name() string         { return "user" }
