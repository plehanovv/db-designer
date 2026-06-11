package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractAttributes(
	document nlp.Document,
	entityMap map[string]*model.Entity,
) {

	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {

		first := tokens[i]
		second := tokens[i+1]
		if isEntity(first) && isAttributeMarker(second.Lemma) {

			entityName := normalizeEntity(first.Lemma)
			entity := entityMap[entityName]
			if entity == nil {
				continue
			}

			for _, token := range collectAttributeTokens(tokens, i+2, first.Sentence) {
				addAttribute(entity, token.Lemma)
			}
		}
	}

	for _, match := range AttributeListMatches(document) {
		entityName := normalizeEntity(match.Owner.Lemma)
		entity := entityMap[entityName]
		if entity == nil {
			continue
		}

		for _, token := range match.Attributes {
			addAttribute(entity, token.Lemma)
		}
	}

	for _, match := range ContainedAttributeMatches(document) {
		entityName := normalizeEntity(match.Owner.Lemma)
		entity := entityMap[entityName]
		if entity == nil {
			continue
		}

		for _, token := range match.Attributes {
			addAttribute(entity, token.Lemma)
		}
	}

	applyAttributeConstraints(document, entityMap)
}

func addAttribute(entity *model.Entity, value string) {
	name := normalizeAttribute(value)
	if name == "" || isIgnoredAttribute(name) || isAttributeMarker(name) {
		return
	}

	for _, attr := range entity.Attributes {
		if attr.Name == name {
			return
		}
	}

	entity.Attributes = append(entity.Attributes, model.Attribute{
		Name:     name,
		Type:     detectAttributeType(name),
		Required: isRequiredAttribute(name),
		Unique:   isUniqueAttribute(name),
	})
}

func applyAttributeConstraints(document nlp.Document, entityMap map[string]*model.Entity) {
	for _, token := range document.Tokens {
		name := normalizeAttribute(token.Lemma)
		if name == "" {
			continue
		}

		required := hasNearbyAttributeConstraint(document.Tokens, token, isRequiredMarker)
		unique := hasNearbyAttributeConstraint(document.Tokens, token, isUniqueMarker)
		if !required && !unique {
			continue
		}

		for _, entity := range entityMap {
			for index := range entity.Attributes {
				if entity.Attributes[index].Name != name {
					continue
				}
				entity.Attributes[index].Required = entity.Attributes[index].Required || required
				entity.Attributes[index].Unique = entity.Attributes[index].Unique || unique
			}
		}
	}
}

func hasNearbyAttributeConstraint(tokens []nlp.Token, attribute nlp.Token, marker func(string) bool) bool {
	for _, token := range tokens {
		if token.Sentence != attribute.Sentence {
			continue
		}
		distance := token.Index - attribute.Index
		if distance < -2 || distance > 2 || distance == 0 {
			continue
		}
		if marker(token.Lemma) {
			return true
		}
	}

	return false
}

func isAttributeMarker(value string) bool {
	markers := map[string]bool{
		"attribute": true, "attributes": true, "field": true, "fields": true,
		"have": true, "has": true,
		"атрибут":  true,
		"атрибуты": true,
		"иметь":    true,
		"имеет":    true,
		"имеют":    true,
		"есть":     true,
		"быть":     true,
		"поле":     true,
		"поля":     true,
	}

	return markers[normalizeWord(value)]
}

func isRequiredAttribute(attribute string) bool {
	required := map[string]bool{
		"id": true, "email": true, "login": true, "name": true,
		"идентификатор": true,
		"почта":         true,
		"логин":         true,
		"название":      true,
		"имя":           true,
	}

	return required[normalizeAttribute(attribute)]
}

func isUniqueAttribute(attribute string) bool {
	unique := map[string]bool{
		"email": true, "login": true, "code": true, "vin": true,
		"почта":    true,
		"логин":    true,
		"код":      true,
		"госномер": true,
		"паспорт":  true,
	}

	return unique[normalizeAttribute(attribute)]
}

func isRequiredMarker(value string) bool {
	markers := map[string]bool{
		"required": true, "mandatory": true,
		"обязательный": true,
		"обязательная": true,
		"обязательное": true,
		"необходимый":  true,
	}

	return markers[normalizeDomainWord(value)]
}

func isUniqueMarker(value string) bool {
	markers := map[string]bool{
		"unique":     true,
		"уникальный": true,
		"уникальная": true,
		"уникальное": true,
		"уникален":   true,
		"уникальна":  true,
	}

	return markers[normalizeDomainWord(value)]
}
