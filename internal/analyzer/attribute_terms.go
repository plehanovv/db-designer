package analyzer

import "db-designer-vkr/internal/nlp"

func AttributeTerms(document nlp.Document) map[string]bool {
	terms := make(map[string]bool)
	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {
		first := tokens[i]
		second := tokens[i+1]

		if !isEntity(first) || !isAttributeMarker(second.Lemma) {
			continue
		}

		for _, token := range collectAttributeTokens(tokens, i+2, first.Sentence) {
			addAttributeTerm(terms, token)
		}
	}

	for _, match := range AttributeListMatches(document) {
		for _, token := range match.Attributes {
			addAttributeTerm(terms, token)
		}
	}

	for _, match := range ContainedAttributeMatches(document) {
		for _, token := range match.Attributes {
			addAttributeTerm(terms, token)
		}
	}

	return terms
}

type AttributeListMatch struct {
	Owner      nlp.Token
	Marker     nlp.Token
	Attributes []nlp.Token
	Rule       string
}

type ContainedAttributeMatch struct {
	Owner      nlp.Token
	Marker     nlp.Token
	Attributes []nlp.Token
	Rule       string
}

func AttributeListMatches(document nlp.Document) []AttributeListMatch {
	tokens := document.Tokens
	var matches []AttributeListMatch

	for i := 0; i < len(tokens)-3; i++ {
		if !isPossessionPrefix(tokens[i].Lemma) || !isEntity(tokens[i+1]) || !isExistenceMarker(tokens[i+2].Lemma) {
			continue
		}

		attributes := collectAttributeTokens(tokens, i+3, tokens[i].Sentence)
		if len(attributes) == 0 {
			continue
		}

		matches = append(matches, AttributeListMatch{
			Owner:      tokens[i+1],
			Marker:     tokens[i+2],
			Attributes: attributes,
			Rule:       "possessive_attribute_list",
		})
	}

	for i := 0; i < len(tokens)-3; i++ {
		if !isEntity(tokens[i]) || !isRelativePronoun(tokens[i+1].Lemma) || !isAttributeMarker(tokens[i+2].Lemma) {
			continue
		}

		attributes := collectAttributeTokens(tokens, i+3, tokens[i].Sentence)
		if len(attributes) == 0 {
			continue
		}

		matches = append(matches, AttributeListMatch{
			Owner:      tokens[i],
			Marker:     tokens[i+2],
			Attributes: attributes,
			Rule:       "relative_attribute_clause",
		})
	}

	for i := 0; i < len(tokens)-4; i++ {
		if !isEntity(tokens[i]) || !isPossessionPrefix(tokens[i+1].Lemma) || !isRelativePronoun(tokens[i+2].Lemma) || !isExistenceMarker(tokens[i+3].Lemma) {
			continue
		}

		attributes := collectAttributeTokens(tokens, i+4, tokens[i].Sentence)
		if len(attributes) == 0 {
			continue
		}

		matches = append(matches, AttributeListMatch{
			Owner:      tokens[i],
			Marker:     tokens[i+3],
			Attributes: attributes,
			Rule:       "relative_possessive_attribute_clause",
		})
	}

	return matches
}

func ContainedAttributeMatches(document nlp.Document) []ContainedAttributeMatch {
	tokens := document.Tokens
	var matches []ContainedAttributeMatch
	emptyAttributeTerms := map[string]bool{}

	for i, token := range tokens {
		if !isContainmentVerb(token.Lemma) {
			continue
		}

		owner, ok := previousEntity(tokens, i, token.Sentence, emptyAttributeTerms)
		if !ok {
			continue
		}

		attributes := collectContainedAttributeTokens(tokens, i+1, token.Sentence)
		if len(attributes) == 0 {
			continue
		}

		matches = append(matches, ContainedAttributeMatch{
			Owner:      owner,
			Marker:     token,
			Attributes: attributes,
			Rule:       "containment_scalar_attribute_list",
		})
	}

	return matches
}

func collectContainedAttributeTokens(tokens []nlp.Token, start int, sentence int) []nlp.Token {
	var attributes []nlp.Token
	for i := start; i < len(tokens); i++ {
		if tokens[i].Sentence != sentence {
			break
		}
		if isRelationPreposition(tokens[i].Lemma) || isRelationQuantifier(tokens[i].Lemma) || isIgnoredAttribute(tokens[i].Lemma) {
			continue
		}
		if isRelationVerb(tokens[i].Lemma) || isAttributeMarker(tokens[i].Lemma) {
			break
		}
		if isKnownScalarAttribute(tokens[i].Lemma) {
			attributes = append(attributes, tokens[i])
		}
	}

	return attributes
}

func collectAttributeTokens(tokens []nlp.Token, start int, sentence int) []nlp.Token {
	var attributes []nlp.Token
	for i := start; i < len(tokens); i++ {
		if tokens[i].Sentence != sentence {
			break
		}
		if isRelationVerb(tokens[i].Lemma) || isAttributeMarker(tokens[i].Lemma) {
			break
		}
		if isRelationQuantifier(tokens[i].Lemma) {
			continue
		}
		if startsWithUpper(tokens[i].Text) && !isKnownScalarAttribute(tokens[i].Lemma) {
			break
		}

		name := normalizeAttribute(tokens[i].Lemma)
		if name == "" || isIgnoredAttribute(name) || isAttributeMarker(name) {
			continue
		}
		if i > start && isRelationQuantifier(tokens[i-1].Lemma) && isEntity(tokens[i]) {
			continue
		}
		if name == "дата" {
			if token, consumed, ok := compoundDateAttribute(tokens, i, sentence); ok {
				attributes = append(attributes, token)
				i += consumed
				continue
			}
		}
		if name == "номер" && i+1 < len(tokens) && tokens[i+1].Sentence == sentence && normalizeAttribute(tokens[i+1].Lemma) == "телефон" {
			continue
		}

		attributes = append(attributes, tokens[i])
	}

	return attributes
}

func compoundDateAttribute(tokens []nlp.Token, index int, sentence int) (nlp.Token, int, bool) {
	if index+1 >= len(tokens) || tokens[index+1].Sentence != sentence {
		return nlp.Token{}, 0, false
	}

	qualifierIndex := index + 1
	qualifier := normalizeAttribute(tokens[qualifierIndex].Lemma)
	if dateModifiers[qualifier] {
		qualifierIndex++
		if qualifierIndex >= len(tokens) || tokens[qualifierIndex].Sentence != sentence {
			return nlp.Token{}, 0, false
		}
		qualifier = normalizeAttribute(tokens[qualifierIndex].Lemma)
	}
	if qualifier == "" || !dateQualifiers[qualifier] {
		return nlp.Token{}, 0, false
	}

	token := tokens[index]
	token.Text = tokens[index].Text + " " + tokens[qualifierIndex].Text
	token.Lemma = "дата_" + qualifier
	return token, qualifierIndex - index, true
}

var dateQualifiers = map[string]bool{
	"заезд":     true,
	"выезд":     true,
	"выдача":    true,
	"возврат":   true,
	"рождение":  true,
	"создание":  true,
	"покупка":   true,
	"начало":    true,
	"окончание": true,
}

var dateModifiers = map[string]bool{
	"первый":    true,
	"первая":    true,
	"первой":    true,
	"последний": true,
	"последняя": true,
	"последней": true,
}

func addAttributeTerm(terms map[string]bool, token nlp.Token) {
	name := normalizeAttribute(token.Lemma)
	if name == "" || isIgnoredAttribute(name) || isAttributeMarker(name) {
		return
	}

	terms[name] = true
}

func isPossessionPrefix(value string) bool {
	prefixes := map[string]bool{
		"у": true,
	}

	return prefixes[normalizeWord(value)]
}

func isExistenceMarker(value string) bool {
	markers := map[string]bool{
		"есть": true,
		"быть": true,
	}

	return markers[normalizeWord(value)]
}

func isRelativePronoun(value string) bool {
	pronouns := map[string]bool{
		"который": true,
		"that":    true,
		"which":   true,
		"who":     true,
	}

	return pronouns[normalizeDomainWord(value)]
}

func isContainmentVerb(value string) bool {
	switch relationType(value) {
	case "contains":
		return true
	default:
		return false
	}
}
