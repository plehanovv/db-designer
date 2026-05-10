package generator

import (
	"db-designer-vkr/internal/model"
	"fmt"
	"strings"
)

func GenerateSQL(entities []model.Entity) string {
	var builder strings.Builder

	for _, entity := range entities {
		tableName := strings.ToLower(entity.Name)

		builder.WriteString(
			fmt.Sprintf(
				"CREATE TABLE %s (\n",
				tableName,
			),
		)

		builder.WriteString("    id SERIAL PRIMARY KEY")

		// атрибуты
		for _, attr := range entity.Attributes {
			builder.WriteString(
				fmt.Sprintf(
					",\n    %s TEXT",
					attr.Name,
				),
			)
		}

		builder.WriteString("\n);\n\n")
	}

	return builder.String()
}
