package analyzer

import (
	"db-designer-vkr/internal/knowledge"
	"db-designer-vkr/internal/model"
)

func ExtractEntities(words []string) map[string]*model.Entity {
	entityMap := make(map[string]*model.Entity)

	for _, word := range words {
		if knowledge.Dictionary[word] == knowledge.EntityType {
			if _, exists := entityMap[word]; !exists {
				entityMap[word] = &model.Entity{
					Name:       word,
					Attributes: []model.Attribute{},
				}
			}
		}
	}

	return entityMap
}
