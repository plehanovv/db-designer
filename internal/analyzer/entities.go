package analyzer

import (
	"strings"

	"db-designer-vkr/internal/model"
)

func ExtractEntities(words []string) map[string]*model.Entity {

	entityMap := make(map[string]*model.Entity)

	for _, w := range words {

		cleanWord := clean(w)

		if isEntity(cleanWord) {

			if _, exists := entityMap[cleanWord]; !exists {

				entityMap[cleanWord] = &model.Entity{
					Name:       cleanWord,
					Attributes: []model.Attribute{},
				}
			}
		}
	}

	return entityMap
}

func isEntity(word string) bool {

	if len(word) == 0 {
		return false
	}

	return strings.ToUpper(string(word[0])) == string(word[0])
}
