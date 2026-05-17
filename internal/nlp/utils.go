package nlp

import "strings"

func clean(word string) string {

	word = strings.ReplaceAll(word, ".", "")
	word = strings.ReplaceAll(word, ",", "")
	word = strings.ReplaceAll(word, ";", "")
	word = strings.ReplaceAll(word, ":", "")

	return strings.TrimSpace(word)
}
