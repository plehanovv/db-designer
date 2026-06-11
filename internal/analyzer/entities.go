package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractEntities(document nlp.Document) map[string]*model.Entity {

	entityMap := make(map[string]*model.Entity)
	attributeTerms := AttributeTerms(document)

	for index, token := range document.Tokens {

		if isEntity(token) {
			if isContextEntityNoise(document.Tokens, index) {
				continue
			}

			if isCompoundDateQualifierToken(document.Tokens, index) {
				continue
			}

			if attributeTerms[normalizeAttribute(token.Lemma)] {
				continue
			}

			entityName := normalizeEntity(token.Lemma)

			if _, exists := entityMap[entityName]; !exists {

				entityMap[entityName] = &model.Entity{
					Name:       entityName,
					Attributes: []model.Attribute{},
				}
			}
		}
	}

	return entityMap
}

func isEntity(token nlp.Token) bool {
	lemma := normalizeWord(token.Lemma)
	if len([]rune(lemma)) < 3 || isIgnoredEntity(lemma) {
		return false
	}

	entityTags := map[string]bool{
		"NOUN":  true,
		"PROPN": true,
	}

	return entityTags[token.Pos]
}

func isContextEntityNoise(tokens []nlp.Token, index int) bool {
	if index <= 0 {
		if index+1 < len(tokens) && isPhoneNumberPhrase(tokens, index) {
			return true
		}
		return false
	}

	value := normalizeWord(tokens[index].Lemma)
	previous := normalizeWord(tokens[index-1].Lemma)

	if isPhoneNumberPhrase(tokens, index) {
		return true
	}

	if previous == "история" && normalizeEntity(value) == "Заказ" {
		return true
	}

	return false
}

func isCompoundDateQualifierToken(tokens []nlp.Token, index int) bool {
	if index <= 0 {
		return false
	}

	token := tokens[index]
	qualifier := normalizeAttribute(token.Lemma)
	if !dateQualifiers[qualifier] {
		return false
	}

	previous := tokens[index-1]
	if previous.Sentence != token.Sentence {
		return false
	}
	if normalizeAttribute(previous.Lemma) == "дата" {
		return true
	}

	if !dateModifiers[normalizeAttribute(previous.Lemma)] || index < 2 {
		return false
	}

	beforeModifier := tokens[index-2]
	return beforeModifier.Sentence == token.Sentence &&
		normalizeAttribute(beforeModifier.Lemma) == "дата"
}

func isPhoneNumberPhrase(tokens []nlp.Token, index int) bool {
	if index+1 >= len(tokens) || tokens[index+1].Sentence != tokens[index].Sentence {
		return false
	}

	return normalizeAttribute(tokens[index].Lemma) == "номер" &&
		normalizeAttribute(tokens[index+1].Lemma) == "телефон"
}
