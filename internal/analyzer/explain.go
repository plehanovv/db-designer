package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

const (
	candidateKindEntity    = "entity"
	candidateKindAttribute = "attribute"
	candidateKindRelation  = "relation"
)

func Explain(document nlp.Document) model.Explanation {
	candidates := make([]model.Candidate, 0)
	candidates = append(candidates, ExplainEntityCandidates(document)...)
	candidates = append(candidates, ExplainAttributeCandidates(document)...)
	candidates = append(candidates, ExplainRelationCandidates(document)...)

	return model.Explanation{
		Candidates: candidates,
	}
}

func ExplainEntityCandidates(document nlp.Document) []model.Candidate {
	var candidates []model.Candidate
	seen := make(map[string]bool)
	attributeTerms := AttributeTerms(document)

	for index, token := range document.Tokens {
		if !isEntity(token) {
			continue
		}

		if isContextEntityNoise(document.Tokens, index) {
			continue
		}

		if isCompoundDateQualifierToken(document.Tokens, index) {
			continue
		}

		if attributeTerms[normalizeAttribute(token.Lemma)] {
			continue
		}

		name := normalizeEntity(token.Lemma)
		if seen[name] {
			continue
		}

		candidates = append(candidates, model.Candidate{
			Kind:       candidateKindEntity,
			Name:       name,
			Rule:       "noun_or_proper_noun",
			SourceText: token.Text,
			Confidence: entityConfidence(token),
			Accepted:   true,
		})
		seen[name] = true
	}

	return candidates
}

func ExplainAttributeCandidates(document nlp.Document) []model.Candidate {
	var candidates []model.Candidate
	seen := make(map[string]bool)
	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {
		first := tokens[i]
		second := tokens[i+1]

		if !isEntity(first) || !isAttributeMarker(second.Lemma) {
			continue
		}

		owner := normalizeEntity(first.Lemma)
		for _, token := range collectAttributeTokens(tokens, i+2, first.Sentence) {
			addAttributeCandidate(&candidates, seen, owner, token, second)
		}
	}

	for _, match := range AttributeListMatches(document) {
		owner := normalizeEntity(match.Owner.Lemma)
		for _, token := range match.Attributes {
			addAttributeCandidateWithRule(&candidates, seen, owner, token, match.Marker, match.Rule)
		}
	}

	for _, match := range ContainedAttributeMatches(document) {
		owner := normalizeEntity(match.Owner.Lemma)
		for _, token := range match.Attributes {
			addAttributeCandidateWithRule(&candidates, seen, owner, token, match.Marker, match.Rule)
		}
	}

	return candidates
}

func ExplainRelationCandidates(document nlp.Document) []model.Candidate {
	var candidates []model.Candidate

	for _, match := range RelationMatches(document) {
		candidates = append(candidates, model.Candidate{
			Kind:       candidateKindRelation,
			Name:       match.Type,
			Owner:      normalizeEntity(match.From.Lemma),
			Target:     normalizeEntity(match.To.Lemma),
			Rule:       match.Rule,
			SourceText: sourceText(match.From, match.Verb, match.To),
			Confidence: relationConfidence(match.Verb.Lemma),
			Accepted:   true,
		})
	}

	return candidates
}

func addAttributeCandidate(
	candidates *[]model.Candidate,
	seen map[string]bool,
	owner string,
	token nlp.Token,
	marker nlp.Token,
) {
	addAttributeCandidateWithRule(candidates, seen, owner, token, marker, "entity_attribute_marker_attribute_list")
}

func addAttributeCandidateWithRule(
	candidates *[]model.Candidate,
	seen map[string]bool,
	owner string,
	token nlp.Token,
	marker nlp.Token,
	rule string,
) {
	name := normalizeAttribute(token.Lemma)
	if name == "" || isIgnoredAttribute(name) || isAttributeMarker(name) {
		return
	}

	key := owner + "." + name
	if seen[key] {
		return
	}

	*candidates = append(*candidates, model.Candidate{
		Kind:       candidateKindAttribute,
		Name:       name,
		Owner:      owner,
		Rule:       rule,
		SourceText: sourceText(marker, token),
		Confidence: attributeConfidence(token),
		Accepted:   true,
	})
	seen[key] = true
}

func entityConfidence(token nlp.Token) float64 {
	if token.Pos == "PROPN" {
		return 0.9
	}
	return 0.78
}

func attributeConfidence(token nlp.Token) float64 {
	if detectAttributeType(token.Lemma) != "TEXT" {
		return 0.86
	}
	return 0.72
}

func relationConfidence(lemma string) float64 {
	if relationType(lemma) == "belongs_to" {
		return 0.88
	}
	return 0.76
}

func sourceText(tokens ...nlp.Token) string {
	value := ""
	for index, token := range tokens {
		if index > 0 {
			value += " "
		}
		value += token.Text
	}
	return value
}
