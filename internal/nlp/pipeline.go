package nlp

type Document struct {
	Tokens   []Token `json:"tokens"`
	Language string  `json:"language,omitempty"`
	Source   string  `json:"source,omitempty"`
}

type Token struct {
	Text       string `json:"text"`
	Lemma      string `json:"lemma"`
	Pos        string `json:"pos"`
	Dependency string `json:"dependency"`
	Head       string `json:"head"`
	Index      int    `json:"index"`
	HeadIndex  int    `json:"headIndex"`
	Sentence   int    `json:"sentence"`
}
