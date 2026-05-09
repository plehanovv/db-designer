package service

import (
	"db-designer-vkr/internal/model"
	"strings"
)

func AnalyzeText(text string) model.AnalyzeResponse {
	words := strings.Fields(text)

	// выделение сущностей
	var entities []model.Entity
	var relations []model.Relation

	seen := make(map[string]bool)

	for _, w := range words {
		cleanWord := clean(w)

		if isEntity(cleanWord) {
			if !seen[cleanWord] {
				entities = append(entities, model.Entity{Name: cleanWord})
				seen[cleanWord] = true
			}
		}
	}

	// выделение связей
	for i := 0; i < len(words)-2; i++ {
		first := clean(words[i])
		second := clean(words[i+1])
		third := clean(words[i+2])

		if isEntity(first) && isVerb(second) && isEntity(third) {
			relations = append(relations, model.Relation{
				From: first,
				To:   third,
				Type: second,
			})
		}
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
