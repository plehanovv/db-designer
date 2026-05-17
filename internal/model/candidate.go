package model

type EntityCandidate struct {
	Name       string
	Confidence float64
	Source     string
}

type AttributeCandidate struct {
	EntityName string
	Name       string
	Type       string
	Confidence float64
}

type RelationCandidate struct {
	From       string
	To         string
	Type       string
	Confidence float64
}
