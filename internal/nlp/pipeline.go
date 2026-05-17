package nlp

import (
	"regexp"
	"strings"
)

type Document struct {
	OriginalText string
	CleanText    string
	Sentences    []Sentence
}

type Sentence struct {
	Text   string
	Tokens []Token
}

type Token struct {
	Value      string
	Normalized string
	Type       string
}
