package slack

import "strings"

func tokenizeSlackCommand(text string) (string, []string) {
	tokens := strings.Split(text, " ")
	return tokens[0], tokens[1:]
}
