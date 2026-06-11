package service

import (
	"fmt"
	"strings"
	"unicode"

	"db-designer-vkr/internal/model"
)

var supportedAttributeTypes = map[string]bool{
	"TEXT":          true,
	"VARCHAR(255)":  true,
	"VARCHAR(20)":   true,
	"INTEGER":       true,
	"NUMERIC(12,2)": true,
	"DATE":          true,
	"TIME":          true,
	"BOOLEAN":       true,
}

var supportedCardinalities = map[string]bool{
	"":             true,
	"one-to-one":   true,
	"one-to-many":  true,
	"many-to-one":  true,
	"many-to-many": true,
	"unspecified":  true,
}

var reservedSQLIdentifiers = map[string]bool{
	"order": true, "select": true, "table": true, "user": true, "where": true,
}

func ValidateModel(database model.Database, entities []model.Entity, relations []model.Relation) []model.Diagnostic {
	var diagnostics []model.Diagnostic

	if strings.TrimSpace(database.Name) == "" {
		diagnostics = append(diagnostics, warning("Database name is empty; SQL generation will skip CREATE DATABASE."))
	}

	if len(entities) == 0 {
		diagnostics = append(diagnostics, errorDiagnostic("No entities were detected; at least one table is required for a database schema."))
		return diagnostics
	}

	entityNames := make(map[string]int)
	entitySQLNames := make(map[string]string)
	for _, entity := range entities {
		name := strings.TrimSpace(entity.Name)
		if name == "" {
			diagnostics = append(diagnostics, errorDiagnostic("Entity with empty name was found."))
			continue
		}

		entityNames[name]++
		if entityNames[name] > 1 {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Duplicate entity name %q was found.", name)))
		}

		sqlName := sqlIdentifier(name)
		if existing, exists := entitySQLNames[sqlName]; exists && existing != name {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Entities %q and %q produce the same SQL table name %q.", existing, name, sqlName)))
		}
		entitySQLNames[sqlName] = name

		if reservedSQLIdentifiers[sqlIdentifierWithoutReservedSuffix(name)] {
			diagnostics = append(diagnostics, warning(fmt.Sprintf("Entity %q uses a reserved SQL identifier; generator will rename it to %q.", name, sqlName)))
		}

		if len(entity.Attributes) == 0 {
			diagnostics = append(diagnostics, warning(fmt.Sprintf("Entity %q has no scalar attributes except the generated primary key.", name)))
		}

		validateAttributes(entity, &diagnostics)
	}

	for _, relation := range relations {
		validateRelation(relation, entityNames, &diagnostics)
	}

	if len(relations) == 0 && len(entities) > 1 {
		diagnostics = append(diagnostics, warning("Several entities were detected, but no relations were found. Add relation phrases or edit relations manually."))
	}

	return diagnostics
}

func validateAttributes(entity model.Entity, diagnostics *[]model.Diagnostic) {
	seenNames := make(map[string]int)
	seenSQLNames := make(map[string]string)

	for _, attribute := range entity.Attributes {
		name := strings.TrimSpace(attribute.Name)
		if name == "" {
			*diagnostics = append(*diagnostics, errorDiagnostic(fmt.Sprintf("Entity %q has an attribute with empty name.", entity.Name)))
			continue
		}

		seenNames[name]++
		if seenNames[name] > 1 {
			*diagnostics = append(*diagnostics, errorDiagnostic(fmt.Sprintf("Entity %q has duplicate attribute %q.", entity.Name, name)))
		}

		sqlName := sqlIdentifier(name)
		if existing, exists := seenSQLNames[sqlName]; exists && existing != name {
			*diagnostics = append(*diagnostics, errorDiagnostic(fmt.Sprintf("Attributes %q and %q of entity %q produce the same SQL column name %q.", existing, name, entity.Name, sqlName)))
		}
		seenSQLNames[sqlName] = name

		if strings.EqualFold(sqlName, "id") {
			*diagnostics = append(*diagnostics, warning(fmt.Sprintf("Entity %q has attribute %q; generated tables already contain primary key id.", entity.Name, name)))
		}

		if !supportedAttributeTypes[strings.ToUpper(strings.TrimSpace(attribute.Type))] {
			*diagnostics = append(*diagnostics, warning(fmt.Sprintf("Attribute %q of entity %q has unsupported type %q.", name, entity.Name, attribute.Type)))
		}
	}
}

func validateRelation(relation model.Relation, entityNames map[string]int, diagnostics *[]model.Diagnostic) {
	from := strings.TrimSpace(relation.From)
	to := strings.TrimSpace(relation.To)

	if from == "" || to == "" {
		*diagnostics = append(*diagnostics, errorDiagnostic("Relation with empty endpoint was found."))
		return
	}

	if entityNames[from] == 0 {
		*diagnostics = append(*diagnostics, errorDiagnostic(fmt.Sprintf("Relation source %q does not match any entity.", from)))
	}
	if entityNames[to] == 0 {
		*diagnostics = append(*diagnostics, errorDiagnostic(fmt.Sprintf("Relation target %q does not match any entity.", to)))
	}
	if from == to {
		*diagnostics = append(*diagnostics, warning(fmt.Sprintf("Relation %q -> %q points to the same entity; check whether this is intentional.", from, to)))
	}
	if !supportedCardinalities[relation.Cardinality] {
		*diagnostics = append(*diagnostics, warning(fmt.Sprintf("Relation %q -> %q has unsupported cardinality %q.", from, to, relation.Cardinality)))
	}
}

func errorDiagnostic(message string) model.Diagnostic {
	return model.Diagnostic{Level: "error", Message: message}
}

func warning(message string) model.Diagnostic {
	return model.Diagnostic{Level: "warning", Message: message}
}

func sqlIdentifier(value string) string {
	result := sqlIdentifierWithoutReservedSuffix(value)
	if reservedSQLIdentifiers[result] {
		return result + "_table"
	}
	return result
}

func sqlIdentifierWithoutReservedSuffix(value string) string {
	value = transliterateSQL(strings.ToLower(value))

	var builder strings.Builder
	previousUnderscore := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			previousUnderscore = false
			continue
		}
		if !previousUnderscore {
			builder.WriteRune('_')
			previousUnderscore = true
		}
	}

	result := strings.Trim(builder.String(), "_")
	if result == "" {
		return "object"
	}
	return result
}

func transliterateSQL(value string) string {
	replacer := strings.NewReplacer(
		"а", "a", "б", "b", "в", "v", "г", "g", "д", "d", "е", "e", "ё", "e",
		"ж", "zh", "з", "z", "и", "i", "й", "y", "к", "k", "л", "l", "м", "m",
		"н", "n", "о", "o", "п", "p", "р", "r", "с", "s", "т", "t", "у", "u",
		"ф", "f", "х", "h", "ц", "c", "ч", "ch", "ш", "sh", "щ", "sch",
		"ъ", "", "ы", "y", "ь", "", "э", "e", "ю", "yu", "я", "ya",
	)

	return replacer.Replace(value)
}
