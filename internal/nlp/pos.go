package nlp

import "strings"

func detectPOS(word string) string {

	lower := strings.ToLower(word)

	verbs := map[string]bool{
		"has":      true,
		"contains": true,
		"belongs":  true,
		"uses":     true,
		"stores":   true,
		"creates":  true,
		"links":    true,
	}

	if verbs[lower] {
		return "VERB"
	}

	if isCapitalized(word) {
		return "NOUN"
	}

	return "UNKNOWN"
}

func isCapitalized(word string) bool {

	if len(word) == 0 {
		return false
	}

	first := string(word[0])

	return first == strings.ToUpper(first)
}
