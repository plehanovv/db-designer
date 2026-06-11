package generator

import (
	"fmt"
	"strings"
	"unicode"

	"db-designer-vkr/internal/model"
)

func GenerateSQL(
	database model.Database,
	entities []model.Entity,
	relations []model.Relation,
) string {

	var builder strings.Builder
	writeDatabaseDDL(&builder, database)

	for _, entity := range entities {
		tableName := toSQLName(entity.Name)

		builder.WriteString(
			fmt.Sprintf(
				"CREATE TABLE %s (\n",
				tableName,
			),
		)

		builder.WriteString("    id SERIAL PRIMARY KEY")

		for _, attr := range entity.Attributes {
			builder.WriteString(
				fmt.Sprintf(
					",\n    %s %s%s",
					toSQLName(attr.Name),
					attr.Type,
					attributeConstraintsSQL(attr),
				),
			)
		}

		for _, foreignKey := range foreignKeysForEntity(entity, relations) {
			builder.WriteString(
				fmt.Sprintf(
					",\n    %s INTEGER,\n    CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(id)",
					foreignKey.columnName,
					tableName,
					foreignKey.targetTableName,
					foreignKey.columnName,
					foreignKey.targetTableName,
				),
			)
		}

		builder.WriteString("\n);\n\n")
	}

	for _, relation := range relations {
		if relation.Cardinality != "many-to-many" {
			continue
		}

		fromTableName := toSQLName(relation.From)
		toTableName := toSQLName(relation.To)
		junctionTableName := fromTableName + "_" + toTableName

		builder.WriteString(
			fmt.Sprintf(
				"CREATE TABLE %s (\n    %s_id INTEGER NOT NULL,\n    %s_id INTEGER NOT NULL,\n    CONSTRAINT pk_%s PRIMARY KEY (%s_id, %s_id),\n    CONSTRAINT fk_%s_%s FOREIGN KEY (%s_id) REFERENCES %s(id),\n    CONSTRAINT fk_%s_%s FOREIGN KEY (%s_id) REFERENCES %s(id)\n);\n\n",
				junctionTableName,
				fromTableName,
				toTableName,
				junctionTableName,
				fromTableName,
				toTableName,
				junctionTableName,
				fromTableName,
				fromTableName,
				fromTableName,
				junctionTableName,
				toTableName,
				toTableName,
				toTableName,
			),
		)
	}

	for _, statement := range indexStatements(entities, relations) {
		builder.WriteString(statement)
		builder.WriteString("\n")
	}

	return builder.String()
}

func writeDatabaseDDL(builder *strings.Builder, database model.Database) {
	name := strings.TrimSpace(database.Name)
	if name == "" {
		return
	}

	databaseName := toSQLName(name)
	builder.WriteString(fmt.Sprintf("CREATE DATABASE %s;\n\n", databaseName))
	builder.WriteString(fmt.Sprintf("-- Connect to %s before running the statements below.\n\n", databaseName))
}

type foreignKey struct {
	columnName      string
	targetTableName string
}

func foreignKeysForEntity(entity model.Entity, relations []model.Relation) []foreignKey {
	var result []foreignKey
	seen := make(map[string]bool)

	for _, relation := range relations {
		holder, target, ok := foreignKeyPlacement(relation)
		if !ok || holder != entity.Name {
			continue
		}

		key := toSQLName(target) + "_id"
		if seen[key] {
			continue
		}

		result = append(result, foreignKey{
			columnName:      key,
			targetTableName: toSQLName(target),
		})
		seen[key] = true
	}

	return result
}

func foreignKeyPlacement(relation model.Relation) (holder string, target string, ok bool) {
	switch relation.Cardinality {
	case "one-to-many":
		return relation.To, relation.From, true
	case "many-to-one":
		return relation.From, relation.To, true
	case "one-to-one", "unspecified", "":
		return relation.From, relation.To, true
	case "many-to-many":
		return "", "", false
	default:
		return relation.From, relation.To, true
	}
}

func attributeConstraintsSQL(attribute model.Attribute) string {
	var constraints []string
	if attribute.Required {
		constraints = append(constraints, "NOT NULL")
	}
	if attribute.Unique {
		constraints = append(constraints, "UNIQUE")
	}
	if len(constraints) == 0 {
		return ""
	}
	return " " + strings.Join(constraints, " ")
}

func indexStatements(entities []model.Entity, relations []model.Relation) []string {
	var statements []string
	seen := make(map[string]bool)

	for _, entity := range entities {
		tableName := toSQLName(entity.Name)
		for _, foreignKey := range foreignKeysForEntity(entity, relations) {
			indexName := "idx_" + tableName + "_" + foreignKey.columnName
			key := indexName + ":" + tableName + ":" + foreignKey.columnName
			if seen[key] {
				continue
			}

			statements = append(statements, fmt.Sprintf("CREATE INDEX %s ON %s(%s);\n", indexName, tableName, foreignKey.columnName))
			seen[key] = true
		}
	}

	for _, relation := range relations {
		if relation.Cardinality != "many-to-many" {
			continue
		}

		fromTableName := toSQLName(relation.From)
		toTableName := toSQLName(relation.To)
		junctionTableName := fromTableName + "_" + toTableName
		indexName := "idx_" + junctionTableName + "_" + toTableName + "_id"
		key := indexName + ":" + junctionTableName
		if seen[key] {
			continue
		}

		statements = append(statements, fmt.Sprintf("CREATE INDEX %s ON %s(%s_id);\n", indexName, junctionTableName, toTableName))
		seen[key] = true
	}

	return statements
}

func toSQLName(value string) string {
	value = transliterate(strings.ToLower(value))

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
	if reservedSQLNames[result] {
		return result + "_table"
	}

	return result
}

var reservedSQLNames = map[string]bool{
	"order": true, "select": true, "table": true, "user": true, "where": true,
}

func transliterate(value string) string {
	replacer := strings.NewReplacer(
		"а", "a", "б", "b", "в", "v", "г", "g", "д", "d", "е", "e", "ё", "e",
		"ж", "zh", "з", "z", "и", "i", "й", "y", "к", "k", "л", "l", "м", "m",
		"н", "n", "о", "o", "п", "p", "р", "r", "с", "s", "т", "t", "у", "u",
		"ф", "f", "х", "h", "ц", "c", "ч", "ch", "ш", "sh", "щ", "sch",
		"ъ", "", "ы", "y", "ь", "", "э", "e", "ю", "yu", "я", "ya",
	)

	return replacer.Replace(value)
}
