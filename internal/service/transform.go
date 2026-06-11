package service

import (
	"fmt"
	"strings"

	"db-designer-vkr/internal/model"
)

func BuildTransformations(database model.Database, entities []model.Entity, relations []model.Relation) []model.TransformationStep {
	var steps []model.TransformationStep

	if strings.TrimSpace(database.Name) != "" {
		steps = append(steps, model.TransformationStep{
			Stage:   "physical",
			Source:  database.Name,
			Target:  sqlIdentifier(database.Name),
			Rule:    "database_name_to_create_database",
			Details: "Subject area name is converted to CREATE DATABASE identifier.",
		})
	}

	for _, entity := range entities {
		tableName := sqlIdentifier(entity.Name)
		steps = append(steps, model.TransformationStep{
			Stage:   "logical",
			Source:  entity.Name,
			Target:  tableName,
			Rule:    "entity_to_table",
			Details: "Entity becomes a relational table with generated primary key id.",
		})

		for _, attribute := range entity.Attributes {
			details := attribute.Type
			if attribute.Required {
				details += ", NOT NULL"
			}
			if attribute.Unique {
				details += ", UNIQUE"
			}
			steps = append(steps, model.TransformationStep{
				Stage:   "logical",
				Source:  entity.Name + "." + attribute.Name,
				Target:  tableName + "." + sqlIdentifier(attribute.Name),
				Rule:    "attribute_to_column",
				Details: details,
			})
		}
	}

	for _, relation := range relations {
		switch relation.Cardinality {
		case "many-to-many":
			fromTable := sqlIdentifier(relation.From)
			toTable := sqlIdentifier(relation.To)
			junction := fromTable + "_" + toTable
			steps = append(steps, model.TransformationStep{
				Stage:   "normalization",
				Source:  relation.From + " -> " + relation.To,
				Target:  junction,
				Rule:    "many_to_many_to_junction_table",
				Details: fmt.Sprintf("Junction table contains %s_id and %s_id foreign keys.", fromTable, toTable),
			})
		default:
			holder, target, ok := relationForeignKeyPlacement(relation)
			if !ok {
				continue
			}
			holderTable := sqlIdentifier(holder)
			targetTable := sqlIdentifier(target)
			column := targetTable + "_id"
			steps = append(steps, model.TransformationStep{
				Stage:   "normalization",
				Source:  relation.From + " -> " + relation.To,
				Target:  holderTable + "." + column,
				Rule:    "relation_to_foreign_key",
				Details: fmt.Sprintf("%s relation places FK on %s table.", relation.Cardinality, holderTable),
			})
		}
	}

	for _, relation := range relations {
		if relation.Cardinality == "many-to-many" {
			fromTable := sqlIdentifier(relation.From)
			toTable := sqlIdentifier(relation.To)
			junction := fromTable + "_" + toTable
			steps = append(steps, model.TransformationStep{
				Stage:   "physical",
				Source:  junction + "." + toTable + "_id",
				Target:  "idx_" + junction + "_" + toTable + "_id",
				Rule:    "foreign_key_index",
				Details: "Index is generated for faster relation traversal.",
			})
			continue
		}

		holder, target, ok := relationForeignKeyPlacement(relation)
		if !ok {
			continue
		}
		holderTable := sqlIdentifier(holder)
		targetTable := sqlIdentifier(target)
		column := targetTable + "_id"
		steps = append(steps, model.TransformationStep{
			Stage:   "physical",
			Source:  holderTable + "." + column,
			Target:  "idx_" + holderTable + "_" + column,
			Rule:    "foreign_key_index",
			Details: "Index is generated for faster joins by foreign key.",
		})
	}

	return steps
}

func relationForeignKeyPlacement(relation model.Relation) (holder string, target string, ok bool) {
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
