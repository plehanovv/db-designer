package nlp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

func ProcessText(text string) Document {
	requestBody := AnalyzeRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Document{}
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Post(
		nlpServiceURL(),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return ProcessTextLocally(text)
	}

	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return ProcessTextLocally(text)
	}

	var document Document

	err = json.NewDecoder(response.Body).Decode(&document)
	if err != nil {
		return ProcessTextLocally(text)
	}

	if len(document.Tokens) == 0 {
		return ProcessTextLocally(text)
	}

	return document
}

func nlpServiceURL() string {
	value := os.Getenv("NLP_SERVICE_URL")
	if value == "" {
		return "http://localhost:8000/analyze"
	}

	return value
}

func ProcessTextLocally(text string) Document {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	tokens := make([]Token, 0, len(words))
	for _, word := range words {
		lemma := strings.ToLower(word)
		tokens = append(tokens, Token{
			Text:  word,
			Lemma: lemma,
			Pos:   detectLocalPOS(lemma),
		})
	}

	return Document{Tokens: tokens}
}

func detectLocalPOS(lemma string) string {
	if localRelationWords[lemma] {
		return "VERB"
	}
	if localNoiseWords[lemma] {
		return "X"
	}
	return "NOUN"
}

var localRelationWords = map[string]bool{
	"have": true, "has": true, "contain": true, "contains": true, "include": true,
	"includes": true, "belong": true, "belongs": true, "connect": true, "connects": true,
	"store": true, "stores": true,
	"\u0438\u043c\u0435\u0435\u0442":                                     true,
	"\u0438\u043c\u0435\u044e\u0442":                                     true,
	"\u0441\u043e\u0434\u0435\u0440\u0436\u0438\u0442":                   true,
	"\u0432\u043a\u043b\u044e\u0447\u0430\u0435\u0442":                   true,
	"\u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0438\u0442": true,
	"\u0441\u0432\u044f\u0437\u0430\u043d":                               true,
	"\u0441\u0432\u044f\u0437\u0430\u043d\u0430":                         true,
	"\u0445\u0440\u0430\u043d\u0438\u0442":                               true,
}

var localNoiseWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "by": true, "each": true,
	"for": true, "from": true, "in": true, "is": true, "of": true, "the": true,
	"to": true, "with": true,
	"\u0432":                               true,
	"\u0438":                               true,
	"\u0438\u043b\u0438":                   true,
	"\u0434\u043b\u044f":                   true,
	"\u043d\u0430":                         true,
	"\u0441":                               true,
	"\u0441\u043e":                         true,
	"\u043f\u043e":                         true,
	"\u0443":                               true,
	"\u043a\u0430\u0436\u0434\u044b\u0439": true,
	"\u043a\u0430\u0436\u0434\u0430\u044f": true,
	"\u043a\u0430\u0436\u0434\u043e\u0435": true,
}
