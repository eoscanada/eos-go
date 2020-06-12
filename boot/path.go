package boot

import (
	"regexp"
)

var weirdities = regexp.MustCompile("[^a-zA-Z0-9]")

func replaceAllWeirdities(input string) string {
	return weirdities.ReplaceAllString(input, "_")
}
