package generator

import (
	"fmt"
	"strings"

	"db-designer-vkr/internal/model"
)

func GenerateSQL(
	entities []model.Entity,
	relations []model.Relation,
) string {

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

		// attributes
		for _, attr := range entity.Attributes {

			builder.WriteString(
				fmt.Sprintf(
					",\n    %s %s",
					strings.ToLower(attr.Name),
					attr.Type,
				),
			)
		}

		// relations -> foreign keys
		for _, relation := range relations {

			if relation.From == entity.Name {

				builder.WriteString(
					fmt.Sprintf(
						",\n    %s_id INTEGER",
						strings.ToLower(relation.To),
					),
				)
			}
		}

		builder.WriteString("\n);\n\n")
	}

	return builder.String()
}
