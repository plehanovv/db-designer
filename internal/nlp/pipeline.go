package nlp

type Document struct {
	OriginalText string
	CleanText    string
	Sentences    []Sentence
}

type Sentence struct {
	Text   string
	Tokens []Token
}
