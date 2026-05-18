package nlp

type Document struct {
	Tokens []Token `json:"tokens"`
}

type Token struct {
	Text       string `json:"text"`
	Lemma      string `json:"lemma"`
	Pos        string `json:"pos"`
	Dependency string `json:"dependency"`
	Head       string `json:"head"`
}
