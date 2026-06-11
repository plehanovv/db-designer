package service

import (
	"sort"
	"strings"

	"db-designer-vkr/internal/analyzer"
	"db-designer-vkr/internal/generator"
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func AnalyzeText(text string) (model.AnalyzeResponse, error) {
	return AnalyzeTextWithDatabase(text, model.Database{})
}

func AnalyzeTextWithDatabase(text string, databaseOverride model.Database) (model.AnalyzeResponse, error) {

	if strings.TrimSpace(text) == "" {
		return model.AnalyzeResponse{}, badInput("domain description is empty")
	}

	if response, parsed, err := tryAnalyzeStructuredInput(text); parsed || err != nil {
		return response, err
	}
	if response, parsed, err := tryAnalyzeCSVInput(text); parsed || err != nil {
		return response, err
	}

	document := nlp.ProcessText(text)
	database := analyzer.ExtractDatabase(text)
	if strings.TrimSpace(databaseOverride.Name) != "" {
		database = databaseOverride
		if strings.TrimSpace(database.Domain) == "" {
			database.Domain = database.Name
		}
	}

	entitiesMap := analyzer.ExtractEntities(document)

	analyzer.ExtractAttributes(document, entitiesMap)

	relations := analyzer.ExtractRelations(document)
	explanation := analyzer.Explain(document)

	var entities []model.Entity

	for _, entity := range entitiesMap {
		entities = append(entities, *entity)
	}

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Name < entities[j].Name
	})

	sql := generator.GenerateSQL(database, entities, relations)
	resultDiagnostics := append(diagnostics(document), ValidateModel(database, entities, relations)...)
	resultDiagnostics = append(resultDiagnostics, ValidateSQL(sql)...)

	response := model.AnalyzeResponse{
		Database:        database,
		Entities:        entities,
		Relations:       relations,
		SQL:             sql,
		Explanation:     explanation,
		Diagnostics:     resultDiagnostics,
		Transformations: BuildTransformations(database, entities, relations),
	}

	return response, nil
}

func GenerateSQL(database model.Database, entities []model.Entity, relations []model.Relation) string {
	return generator.GenerateSQL(database, entities, relations)
}

func GenerateSQLWithDiagnostics(database model.Database, entities []model.Entity, relations []model.Relation) (string, []model.Diagnostic) {
	sql := generator.GenerateSQL(database, entities, relations)
	diagnostics := append(ValidateModel(database, entities, relations), ValidateSQL(sql)...)
	return sql, diagnostics
}

func GenerateSQLWithMetadata(database model.Database, entities []model.Entity, relations []model.Relation) (string, []model.Diagnostic, []model.TransformationStep) {
	sql := generator.GenerateSQL(database, entities, relations)
	diagnostics := append(ValidateModel(database, entities, relations), ValidateSQL(sql)...)
	return sql, diagnostics, BuildTransformations(database, entities, relations)
}

func diagnostics(document nlp.Document) []model.Diagnostic {
	var result []model.Diagnostic

	switch document.Source {
	case "spacy":
		result = append(result, model.Diagnostic{
			Level:   "info",
			Message: "NLP analysis uses spaCy tokens, lemmas, POS tags and dependency metadata.",
		})
	default:
		result = append(result, model.Diagnostic{
			Level:   "warning",
			Message: "NLP service is unavailable or returned weak tags; local rule-based fallback is used.",
		})
	}

	if len(document.Tokens) == 0 {
		result = append(result, model.Diagnostic{
			Level:   "warning",
			Message: "No NLP tokens were produced from the input text.",
		})
	}

	return result
}
