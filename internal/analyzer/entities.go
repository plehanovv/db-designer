package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractEntities(document nlp.Document) map[string]*model.Entity {

	entityMap := make(map[string]*model.Entity)

	for _, token := range document.Tokens {

		if isEntity(token) {

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
