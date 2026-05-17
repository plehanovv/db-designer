package nlp

import (
	"regexp"
	"strings"
)

func ProcessText(text string) Document {

	clean := normalizeText(text)

	sentences := splitSentences(clean)

	var result []Sentence

	for _, s := range sentences {

		tokens := tokenize(s)

		result = append(result, Sentence{
			Text:   s,
			Tokens: tokens,
		})
	}

	return Document{
		OriginalText: text,
		CleanText:    clean,
		Sentences:    result,
	}
}

func normalizeText(text string) string {

	text = strings.TrimSpace(text)

	re := regexp.MustCompile(`[^\w\s\.]`)

	text = re.ReplaceAllString(text, "")

	return text
}

func splitSentences(text string) []string {

	raw := strings.Split(text, ".")

	var result []string

	for _, s := range raw {

		s = strings.TrimSpace(s)

		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

func tokenize(sentence string) []Token {

	words := strings.Fields(sentence)

	var tokens []Token

	for _, w := range words {

		tokens = append(tokens, Token{
			Value:      w,
			Normalized: strings.ToLower(w),
			Type:       detectTokenType(w),
		})
	}

	return tokens
}
