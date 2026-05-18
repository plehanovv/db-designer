package analyzer

import (
	"strings"

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

	entityTags := map[string]bool{
		"NOUN":  true,
		"PROPN": true,
	}

	return entityTags[token.Pos]
}

func normalizeEntity(value string) string {

	value = strings.TrimSpace(value)

	if value == "" {
		return value
	}

	first := strings.ToUpper(string(value[0]))

	return first + value[1:]
}
