package analyzer

import "db-designer-vkr/internal/nlp"

type RelationMatch struct {
	Verb        nlp.Token
	From        nlp.Token
	To          nlp.Token
	Type        string
	Cardinality string
	Rule        string
}

func RelationMatches(document nlp.Document) []RelationMatch {
	var matches []RelationMatch
	seen := make(map[string]bool)
	attributeTerms := AttributeTerms(document)
	tokens := document.Tokens

	for _, match := range ownershipRelationMatches(tokens, attributeTerms) {
		addRelationMatch(&matches, seen, match)
	}

	for i, token := range tokens {
		if !isRelationVerb(token.Lemma) {
			continue
		}

		fromToken, toToken, rule, ok := dependencyRelationEndpoints(tokens, token, attributeTerms)
		if !ok {
			fromToken, ok = previousEntity(tokens, i, token.Sentence, attributeTerms)
			if !ok {
				continue
			}

			toToken, ok = nextEntity(tokens, i+1, token.Sentence, attributeTerms)
			if !ok {
				continue
			}

			rule = "linear_relation_marker"
		}

		addRelationMatch(&matches, seen, RelationMatch{
			Verb:        token,
			From:        fromToken,
			To:          toToken,
			Type:        relationType(token.Lemma),
			Cardinality: relationCardinality(token.Lemma),
			Rule:        rule,
		})
	}

	return matches
}

func addRelationMatch(matches *[]RelationMatch, seen map[string]bool, match RelationMatch) {
	from := normalizeEntity(match.From.Lemma)
	to := normalizeEntity(match.To.Lemma)
	if from == to {
		return
	}

	key := from + "->" + to + ":" + match.Type
	if seen[key] {
		return
	}

	*matches = append(*matches, match)
	seen[key] = true
}

func ownershipRelationMatches(tokens []nlp.Token, attributeTerms map[string]bool) []RelationMatch {
	var matches []RelationMatch

	for i, token := range tokens {
		if !isAttributeMarker(token.Lemma) {
			continue
		}

		fromToken, fromOK := previousEntity(tokens, i, token.Sentence, attributeTerms)
		toToken, toOK := nextEntity(tokens, i+1, token.Sentence, attributeTerms)
		if !fromOK || !toOK {
			continue
		}

		cardinality, ok := ownershipCardinality(tokens, fromToken, token, toToken)
		if !ok {
			continue
		}

		matches = append(matches, RelationMatch{
			Verb:        token,
			From:        fromToken,
			To:          toToken,
			Type:        "has",
			Cardinality: cardinality,
			Rule:        "quantified_possession_relation",
		})
	}

	return matches
}

func ownershipCardinality(tokens []nlp.Token, from nlp.Token, marker nlp.Token, to nlp.Token) (string, bool) {
	sourceOne := hasQuantifierBefore(tokens, from, []string{"one", "один", "одна", "одно", "each", "каждый", "каждая"})
	targetMany := hasQuantifierBetween(tokens, marker.Index, to.Index, []string{"many", "multiple", "several", "много", "несколько"})

	if targetMany {
		return "one-to-many", true
	}
	if sourceOne {
		return "one-to-many", true
	}

	return "", false
}

func hasQuantifierBefore(tokens []nlp.Token, target nlp.Token, values []string) bool {
	for _, token := range tokens {
		if token.Sentence != target.Sentence || token.Index >= target.Index {
			continue
		}
		if containsNormalized(values, token.Lemma) {
			return true
		}
	}

	return false
}

func hasQuantifierBetween(tokens []nlp.Token, start int, end int, values []string) bool {
	for _, token := range tokens {
		if token.Index <= start || token.Index >= end {
			continue
		}
		if containsNormalized(values, token.Lemma) {
			return true
		}
	}

	return false
}

func containsNormalized(values []string, value string) bool {
	value = normalizeWord(value)
	for _, candidate := range values {
		if normalizeWord(candidate) == value {
			return true
		}
	}

	return false
}

func dependencyRelationEndpoints(
	tokens []nlp.Token,
	verb nlp.Token,
	attributeTerms map[string]bool,
) (nlp.Token, nlp.Token, string, bool) {
	if verb.Index == 0 && verb.HeadIndex == 0 && verb.Dependency == "" {
		return nlp.Token{}, nlp.Token{}, "", false
	}

	subject, subjectOK := dependencySubject(tokens, verb, attributeTerms)
	object, objectOK := dependencyObject(tokens, verb, attributeTerms)
	if !subjectOK || !objectOK {
		return nlp.Token{}, nlp.Token{}, "", false
	}

	return subject, object, "dependency_subject_predicate_object", true
}

func dependencySubject(tokens []nlp.Token, verb nlp.Token, attributeTerms map[string]bool) (nlp.Token, bool) {
	for _, token := range tokens {
		if token.Sentence != verb.Sentence || token.HeadIndex != verb.Index {
			continue
		}

		if isSubjectDependency(token.Dependency) && isRelationEntity(token, attributeTerms) {
			return token, true
		}
	}

	return nlp.Token{}, false
}

func dependencyObject(tokens []nlp.Token, verb nlp.Token, attributeTerms map[string]bool) (nlp.Token, bool) {
	for _, token := range tokens {
		if token.Sentence != verb.Sentence || !isRelationEntity(token, attributeTerms) {
			continue
		}

		if token.HeadIndex == verb.Index && isObjectDependency(token.Dependency) {
			return token, true
		}

		if isPrepositionalObject(token, tokens, verb) {
			return token, true
		}
	}

	return nlp.Token{}, false
}

func isPrepositionalObject(token nlp.Token, tokens []nlp.Token, verb nlp.Token) bool {
	if !isObjectDependency(token.Dependency) {
		return false
	}

	head, ok := tokenByIndex(tokens, token.HeadIndex)
	if !ok || head.Sentence != verb.Sentence {
		return false
	}

	if head.HeadIndex == verb.Index && isRelationPreposition(head.Lemma) {
		return true
	}

	return token.HeadIndex == verb.Index
}

func tokenByIndex(tokens []nlp.Token, index int) (nlp.Token, bool) {
	for _, token := range tokens {
		if token.Index == index {
			return token, true
		}
	}

	return nlp.Token{}, false
}

func isRelationEntity(token nlp.Token, attributeTerms map[string]bool) bool {
	return isEntity(token) && !attributeTerms[normalizeAttribute(token.Lemma)]
}

func isSubjectDependency(value string) bool {
	dependencies := map[string]bool{
		"nsubj": true, "nsubj:pass": true, "nsubjpass": true,
	}

	return dependencies[normalizeWord(value)]
}

func isObjectDependency(value string) bool {
	dependencies := map[string]bool{
		"obj": true, "dobj": true, "iobj": true, "pobj": true,
		"obl": true, "nmod": true,
	}

	return dependencies[normalizeWord(value)]
}
