package model

type Database struct {
	Name   string `json:"name"`
	Domain string `json:"domain,omitempty"`
}

type Entity struct {
	Name       string      `json:"name"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Unique   bool   `json:"unique,omitempty"`
}

type Relation struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`
	Cardinality string `json:"cardinality"`
}

type Candidate struct {
	Kind       string  `json:"kind"`
	Name       string  `json:"name"`
	Owner      string  `json:"owner,omitempty"`
	Target     string  `json:"target,omitempty"`
	Rule       string  `json:"rule"`
	SourceText string  `json:"sourceText"`
	Confidence float64 `json:"confidence"`
	Accepted   bool    `json:"accepted"`
}

type Explanation struct {
	Candidates []Candidate `json:"candidates"`
}

type Diagnostic struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type TransformationStep struct {
	Stage   string `json:"stage"`
	Source  string `json:"source"`
	Target  string `json:"target"`
	Rule    string `json:"rule"`
	Details string `json:"details,omitempty"`
}

type AnalyzeResponse struct {
	Database        Database             `json:"database"`
	Entities        []Entity             `json:"entities"`
	Relations       []Relation           `json:"relations"`
	SQL             string               `json:"sql"`
	Explanation     Explanation          `json:"explanation"`
	Diagnostics     []Diagnostic         `json:"diagnostics"`
	Transformations []TransformationStep `json:"transformations"`
	StorageKey      string               `json:"storageKey,omitempty"`
}
