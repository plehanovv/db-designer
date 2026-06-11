package service

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"db-designer-vkr/internal/generator"
	"db-designer-vkr/internal/model"
)

type structuredAnalyzeInput struct {
	Database  model.Database   `json:"database"`
	Entities  []model.Entity   `json:"entities"`
	Relations []model.Relation `json:"relations"`
}

func tryAnalyzeStructuredInput(text string) (model.AnalyzeResponse, bool, error) {
	if !looksLikeJSON(text) {
		return model.AnalyzeResponse{}, false, nil
	}

	var input structuredAnalyzeInput
	if err := json.Unmarshal([]byte(text), &input); err != nil {
		return model.AnalyzeResponse{}, true, badInput("invalid structured JSON model")
	}
	if len(input.Entities) == 0 {
		return model.AnalyzeResponse{}, true, badInput("structured JSON model must contain entities")
	}

	sql := generator.GenerateSQL(input.Database, input.Entities, input.Relations)
	diagnostics := []model.Diagnostic{
		{Level: "info", Message: "Structured JSON input was parsed directly; NLP analysis was skipped."},
	}
	diagnostics = append(diagnostics, ValidateModel(input.Database, input.Entities, input.Relations)...)
	diagnostics = append(diagnostics, ValidateSQL(sql)...)

	return model.AnalyzeResponse{
		Database:        input.Database,
		Entities:        input.Entities,
		Relations:       input.Relations,
		SQL:             sql,
		Explanation:     structuredExplanation(input.Entities, input.Relations, "structured_json"),
		Diagnostics:     diagnostics,
		Transformations: BuildTransformations(input.Database, input.Entities, input.Relations),
	}, true, nil
}

func tryAnalyzeCSVInput(text string) (model.AnalyzeResponse, bool, error) {
	if !looksLikeCSV(text) {
		return model.AnalyzeResponse{}, false, nil
	}

	reader := csv.NewReader(strings.NewReader(text))
	reader.TrimLeadingSpace = true
	rows, err := reader.ReadAll()
	if err != nil {
		return model.AnalyzeResponse{}, true, badInput("invalid structured CSV model")
	}
	if len(rows) < 2 {
		return model.AnalyzeResponse{}, true, badInput("structured CSV model must contain a header and at least one data row")
	}

	header := csvHeaderIndex(rows[0])
	entitiesByName := make(map[string]*model.Entity)
	var order []string
	var relations []model.Relation
	database := model.Database{Name: firstCSVDatabaseName(header, rows[1:])}

	for _, row := range rows[1:] {
		kind := strings.ToLower(csvCell(header, row, "kind"))
		entityName := csvCell(header, row, "entity")
		if entityName != "" {
			if _, exists := entitiesByName[entityName]; !exists {
				entitiesByName[entityName] = &model.Entity{Name: entityName, Attributes: []model.Attribute{}}
				order = append(order, entityName)
			}
		}

		switch kind {
		case "relation":
			relations = append(relations, model.Relation{
				From:        csvCell(header, row, "from"),
				To:          csvCell(header, row, "to"),
				Type:        defaultString(csvCell(header, row, "relation_type"), "associated_with"),
				Cardinality: defaultString(csvCell(header, row, "cardinality"), "unspecified"),
			})
		default:
			attributeName := csvCell(header, row, "attribute")
			if entityName == "" || attributeName == "" {
				continue
			}
			entity := entitiesByName[entityName]
			entity.Attributes = append(entity.Attributes, model.Attribute{
				Name:     attributeName,
				Type:     defaultString(csvCell(header, row, "type"), "TEXT"),
				Required: parseBool(csvCell(header, row, "required")),
				Unique:   parseBool(csvCell(header, row, "unique")),
			})
		}
	}

	var entities []model.Entity
	for _, name := range order {
		entities = append(entities, *entitiesByName[name])
	}
	if len(entities) == 0 {
		return model.AnalyzeResponse{}, true, badInput("structured CSV model must contain entities")
	}

	sql := generator.GenerateSQL(database, entities, relations)
	diagnostics := []model.Diagnostic{{Level: "info", Message: "Structured CSV input was parsed directly; NLP analysis was skipped."}}
	diagnostics = append(diagnostics, ValidateModel(database, entities, relations)...)
	diagnostics = append(diagnostics, ValidateSQL(sql)...)

	return model.AnalyzeResponse{
		Database:        database,
		Entities:        entities,
		Relations:       relations,
		SQL:             sql,
		Explanation:     structuredExplanation(entities, relations, "structured_csv"),
		Diagnostics:     diagnostics,
		Transformations: BuildTransformations(database, entities, relations),
	}, true, nil
}

func looksLikeJSON(text string) bool {
	return LooksLikeJSON(text)
}

func LooksLikeJSON(text string) bool {
	text = strings.TrimSpace(text)
	return strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}")
}

func looksLikeCSV(text string) bool {
	return LooksLikeCSV(text)
}

func LooksLikeCSV(text string) bool {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "{") || !strings.Contains(text, ",") {
		return false
	}
	reader := csv.NewReader(strings.NewReader(text))
	header, err := reader.Read()
	if err != nil && err != io.EOF {
		return false
	}
	index := csvHeaderIndex(header)
	_, hasEntity := index["entity"]
	_, hasKind := index["kind"]
	return hasEntity || hasKind
}

func csvHeaderIndex(header []string) map[string]int {
	index := make(map[string]int)
	for i, value := range header {
		index[strings.ToLower(strings.TrimSpace(value))] = i
	}
	return index
}

func csvCell(header map[string]int, row []string, name string) string {
	index, exists := header[name]
	if !exists || index < 0 || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func firstCSVDatabaseName(header map[string]int, rows [][]string) string {
	for _, row := range rows {
		if value := csvCell(header, row, "database"); value != "" {
			return value
		}
	}
	return ""
}

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "y", "да":
		return true
	default:
		return false
	}
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func structuredExplanation(entities []model.Entity, relations []model.Relation, source string) model.Explanation {
	var candidates []model.Candidate

	for _, entity := range entities {
		candidates = append(candidates, model.Candidate{
			Kind:       "entity",
			Name:       entity.Name,
			Rule:       source + "_entity",
			SourceText: entity.Name,
			Confidence: 1,
			Accepted:   true,
		})

		for _, attribute := range entity.Attributes {
			candidates = append(candidates, model.Candidate{
				Kind:       "attribute",
				Name:       attribute.Name,
				Owner:      entity.Name,
				Rule:       source + "_attribute",
				SourceText: attribute.Name,
				Confidence: 1,
				Accepted:   true,
			})
		}
	}

	for _, relation := range relations {
		candidates = append(candidates, model.Candidate{
			Kind:       "relation",
			Name:       relation.Type,
			Owner:      relation.From,
			Target:     relation.To,
			Rule:       source + "_relation",
			SourceText: relation.From + " " + relation.Type + " " + relation.To,
			Confidence: 1,
			Accepted:   true,
		})
	}

	return model.Explanation{Candidates: candidates}
}
