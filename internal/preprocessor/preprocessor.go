package preprocessor

import (
	"regexp"
	"strings"
)

var stopWords = map[string]bool{
	"the": true,
	"and": true,
	"or":  true,
	"is":  true,
	"a":   true,
	"an":  true,
	"of":  true,
	"in":  true,
	"on":  true,
	"for": true,
	"to":  true,
}

func Process(text string) []string {
	text = normalize(text)

	tokens := tokenize(text)

	tokens = removeStopWords(tokens)

	return tokens
}

func normalize(text string) string {
	text = strings.ToLower(text)

	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)

	text = reg.ReplaceAllString(text, "")

	return text
}

func tokenize(text string) []string {
	return strings.Fields(text)
}

func removeStopWords(words []string) []string {
	var filtered []string

	for _, word := range words {

		if !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}
