package analyzer

import "strings"

func clean(word string) string {

	word = strings.ReplaceAll(word, ".", "")
	word = strings.ReplaceAll(word, ",", "")

	return word
}
