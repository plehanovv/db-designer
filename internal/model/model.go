package model

type Entity struct {
	Name string `json:"name"`
}

type Relation struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type AnalyzeResponse struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
}
