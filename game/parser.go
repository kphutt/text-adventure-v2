package game

import "strings"

// ParseInput splits a string into a verb and a noun.
func ParseInput(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}
	verb := parts[0]
	noun := ""
	if len(parts) > 1 {
		noun = strings.Join(parts[1:], " ")
	}
	return verb, noun
}
