package service

import (
	"errors"
	"sort"
	"strings"

	"db-designer-vkr/internal/analyzer"
	"db-designer-vkr/internal/generator"
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func AnalyzeText(text string) (model.AnalyzeResponse, error) {

	if strings.TrimSpace(text) == "" {
		return model.AnalyzeResponse{}, errors.New("empty text")
	}

	document := nlp.ProcessText(text)

	entitiesMap := analyzer.ExtractEntities(document)

	analyzer.ExtractAttributes(document, entitiesMap)

	relations := analyzer.ExtractRelations(document)

	var entities []model.Entity

	for _, entity := range entitiesMap {
		entities = append(entities, *entity)
	}

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Name < entities[j].Name
	})

	sql := generator.GenerateSQL(entities, relations)

	response := model.AnalyzeResponse{
		Entities:  entities,
		Relations: relations,
		SQL:       sql,
	}

	return response, nil
}
