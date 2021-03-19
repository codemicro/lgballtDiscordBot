package bios

import (
	"errors"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"regexp"
	"strings"
)

var pluralkitMemberIdRegexp = regexp.MustCompile(`(?m)^[a-zA-Z]{5}$`)

type pluralkitMemberIdType struct{}

func (pluralkitMemberIdType) Parse(content *string) (interface{}, error) {

	a, b := common.TakeFirstPart(*content)

	if pluralkitMemberIdRegexp.MatchString(a) {
		*content = b
		return strings.ToLower(a), nil
	}

	return nil, errors.New("not a valid PluralKit member ID")

}
func (pluralkitMemberIdType) Help(_ string) string {
	return "A PluralKit member ID, for example `abcde`"
}
func (pluralkitMemberIdType) Name() string { return "pkMemberId" }

type bioFieldType struct{}

func (bioFieldType) Parse(content *string) (interface{}, error) {

	a, b := common.TakeFirstPart(*content)

	var found bool
	var properFieldName string
	for _, f := range config.BioFields {
		if strings.EqualFold(f, a) {
			found = true
			properFieldName = f
			break
		}
	}

	if found {
		*content = b
		return properFieldName, nil
	}

	return nil, errors.New("invalid field name: options are: " + strings.Join(config.BioFields, ", "))

}
func (bioFieldType) Help(_ string) string { return "A bio field name, eg `pronouns` or `gender`" }
func (bioFieldType) Name() string         { return "bioField" }
