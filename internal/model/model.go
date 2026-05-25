package model

type Entity struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Relation struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`
	Cardinality string `json:"cardinality"`
}

type AnalyzeResponse struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
	SQL       string     `json:"sql"`
}
