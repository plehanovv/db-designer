package generator

import (
	"fmt"
	"strings"
	"unicode"

	"db-designer-vkr/internal/model"
)

func GenerateSQL(
	entities []model.Entity,
	relations []model.Relation,
) string {

	var builder strings.Builder

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
					requiredSQL(attr.Required),
				),
			)
		}

		for _, relation := range relations {
			if relation.From == entity.Name {
				columnName := toSQLName(relation.To) + "_id"
				targetTableName := toSQLName(relation.To)

				builder.WriteString(
					fmt.Sprintf(
						",\n    %s INTEGER,\n    CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(id)",
						columnName,
						tableName,
						targetTableName,
						columnName,
						targetTableName,
					),
				)
			}
		}

		builder.WriteString("\n);\n\n")
	}

	return builder.String()
}

func requiredSQL(required bool) string {
	if required {
		return " NOT NULL"
	}
	return ""
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
		"\u0430", "a", "\u0431", "b", "\u0432", "v", "\u0433", "g", "\u0434", "d", "\u0435", "e", "\u0451", "e",
		"\u0436", "zh", "\u0437", "z", "\u0438", "i", "\u0439", "y", "\u043a", "k", "\u043b", "l", "\u043c", "m",
		"\u043d", "n", "\u043e", "o", "\u043f", "p", "\u0440", "r", "\u0441", "s", "\u0442", "t", "\u0443", "u",
		"\u0444", "f", "\u0445", "h", "\u0446", "c", "\u0447", "ch", "\u0448", "sh", "\u0449", "sch",
		"\u044a", "", "\u044b", "y", "\u044c", "", "\u044d", "e", "\u044e", "yu", "\u044f", "ya",
	)

	return replacer.Replace(value)
}
