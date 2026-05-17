package nlp

import "strings"

func detectTokenType(word string) string {

	if isVerb(word) {
		return "VERB"
	}

	if isPotentialEntity(word) {
		return "ENTITY"
	}

	return "WORD"
}

func isPotentialEntity(word string) bool {

	if len(word) == 0 {
		return false
	}

	first := string(word[0])

	return strings.ToUpper(first) == first
}
