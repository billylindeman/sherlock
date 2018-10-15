package sherlock

import (
	"regexp"
	"strings"
)

// token holds information about a Token found in a document
type token struct {
	value    string
	position int
	fieldID  int
	weight   int
}

type analysis struct {
	tokens []token
}

var reg, _ = regexp.Compile("[^a-zA-Z0-9\\s]+")

func normalize(in string) string {
	lower := strings.ToLower(in)
	charRemoved := reg.ReplaceAllString(lower, "")
	return charRemoved
}
