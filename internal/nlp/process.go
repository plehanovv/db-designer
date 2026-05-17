package nlp

import "strings"

func ProcessText(text string) Document {

	rawWords := strings.Fields(text)

	var tokens []Token

	for _, word := range rawWords {

		cleaned := clean(word)

		token := Token{
			Value: cleaned,
			Lemma: strings.ToLower(cleaned),
			POS:   detectPOS(cleaned),
		}

		tokens = append(tokens, token)
	}

	return Document{
		Text:   text,
		Tokens: tokens,
	}
}
