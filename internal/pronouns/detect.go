package pronouns

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

var pronounRegexp = regexp.MustCompile(`(?mi)(?:(?:[A-Za-z]+/[A-Za-z]+)|any pronouns|no pronouns)`)

type PronounRole struct {
	Name   string
	RoleID string
}

// FilterRoleList filters a slice of Discord roles depending on if they match an expected pronoun format.
func FilterRoleList(roles []*discordgo.Role) (o []PronounRole) {
	for _, role := range roles {
		if pronounRegexp.MatchString(role.Name) {
			o = append(o, PronounRole{
				Name:   role.Name,
				RoleID: role.ID,
			})
		}
	}
	return
}

// FindPronounsInString searches for and locates pronouns in the form he/him, he/they or him/them based of the provided
// set of pronoun roles.
func FindPronounsInString(in string, possiblePronouns []PronounRole) (o []PronounRole) {

	matches := pronounRegexp.FindAllString(in, -1)
	if matches == nil {
		return
	}

	addedPronouns := make(map[string]struct{})

	for _, match := range matches {

		splitMatch := strings.Split(match, "/")

		for _, matchPart := range splitMatch {

			for _, pronoun := range possiblePronouns {

				if _, found := addedPronouns[pronoun.Name]; found {
					continue
				}

				splitPronoun := strings.Split(pronoun.Name, "/")

				for _, pronounPart := range splitPronoun {
					if strings.EqualFold(pronounPart, matchPart) {
						o = append(o, pronoun)
						addedPronouns[pronoun.Name] = struct{}{}
						break
					}
				}

			}

		}

	}

	return
}
