package service

import (
	"db-designer-vkr/internal/model"
	"strings"
)

func AnalyzeText(text string) model.AnalyzeResponse {
	words := strings.Fields(text)

	entityMap := make(map[string]*model.Entity)
	var relations []model.Relation

	// выделение сущностей
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

	// выделение связей и атрибутов
	for i := 0; i < len(words)-2; i++ {
		first := clean(words[i])
		second := clean(words[i+1])
		third := clean(words[i+2])

		// Entity + verb + entity
		if isEntity(first) && isVerb(second) && isEntity(third) {
			relations = append(relations, model.Relation{
				From: first,
				To:   third,
				Type: second,
			})
		}

		if isEntity(first) && second == "has" && !isEntity(third) {
			entity := entityMap[first]

			entity.Attributes = append(entity.Attributes, model.Attribute{
				Name: third,
				Type: "string",
			})
		}
	}

	var entities []model.Entity

	for _, entity := range entityMap {
		entities = append(entities, *entity)
	}

	return model.AnalyzeResponse{
		Entities:  entities,
		Relations: relations,
	}
}

func isEntity(word string) bool {
	if len(word) == 0 {
		return false
	}

	return strings.ToUpper(string(word[0])) == string(word[0])
}

func isVerb(word string) bool {
	verbs := []string{
		"has",
		"contains",
		"enrolls",
		"teaches",
		"belongs",
		"uses",
		"manages",
		"creates",
	}

	for _, verb := range verbs {
		if word == verb {
			return true
		}
	}

	return false
}

func clean(word string) string {
	word = strings.ReplaceAll(word, ".", "")
	word = strings.ReplaceAll(word, ",", "")

	return word
}
