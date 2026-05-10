package service

import (
	"strings"

	"db-designer-vkr/internal/analyzer"
	"db-designer-vkr/internal/generator"
	"db-designer-vkr/internal/model"
)

func AnalyzeText(text string) model.AnalyzeResponse {

	words := strings.Fields(text)

	entitiesMap := analyzer.ExtractEntities(words)

	analyzer.ExtractAttributes(words, entitiesMap)

	relations := analyzer.ExtractRelations(words)

	var entities []model.Entity

	for _, entity := range entitiesMap {
		entities = append(entities, *entity)
	}

	sql := generator.GenerateSQL(entities)

	return model.AnalyzeResponse{
		Entities:  entities,
		Relations: relations,
		SQL:       sql,
	}
}
