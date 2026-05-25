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
		third := tokens[i+2]

		if isEntity(first) && isAttributeMarker(second.Lemma) {

			entityName := normalizeEntity(first.Lemma)
			entity := entityMap[entityName]
			if entity == nil {
				continue
			}

			addAttribute(entity, third.Lemma)

			for j := i + 3; j < len(tokens); j++ {
				if isRelationVerb(tokens[j].Lemma) || isAttributeMarker(tokens[j].Lemma) {
					break
				}

				if startsWithUpper(tokens[j].Text) {
					break
				}

				addAttribute(entity, tokens[j].Lemma)
			}
		}
	}
}

func addAttribute(entity *model.Entity, value string) {
	name := normalizeAttribute(value)
	if name == "" || isAttributeMarker(name) {
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
	})
}

func isAttributeMarker(value string) bool {
	markers := map[string]bool{
		"have": true, "has": true, "include": true, "includes": true,
		"contain": true, "contains": true,
		"\u0438\u043c\u0435\u0442\u044c":                         true,
		"\u0438\u043c\u0435\u0435\u0442":                         true,
		"\u0438\u043c\u0435\u044e\u0442":                         true,
		"\u0441\u043e\u0434\u0435\u0440\u0436\u0430\u0442\u044c": true,
		"\u0441\u043e\u0434\u0435\u0440\u0436\u0438\u0442":       true,
		"\u0432\u043a\u043b\u044e\u0447\u0430\u0442\u044c":       true,
		"\u0432\u043a\u043b\u044e\u0447\u0430\u0435\u0442":       true,
	}

	return markers[normalizeWord(value)]
}

func isRequiredAttribute(attribute string) bool {
	required := map[string]bool{
		"id": true, "email": true, "login": true, "name": true,
		"\u0438\u0434\u0435\u043d\u0442\u0438\u0444\u0438\u043a\u0430\u0442\u043e\u0440": true,
		"\u043f\u043e\u0447\u0442\u0430":                                                 true,
		"\u043b\u043e\u0433\u0438\u043d":                                                 true,
		"\u043d\u0430\u0437\u0432\u0430\u043d\u0438\u0435":                               true,
		"\u0438\u043c\u044f":                                                             true,
	}

	return required[normalizeAttribute(attribute)]
}
