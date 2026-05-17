package nlp

import "strings"

var verbs = map[string]bool{
	"has":        true,
	"contains":   true,
	"belongs":    true,
	"includes":   true,
	"uses":       true,
	"creates":    true,
	"stores":     true,
	"manages":    true,
	"references": true,
}

func isVerb(word string) bool {

	_, exists := verbs[strings.ToLower(word)]

	return exists
}
