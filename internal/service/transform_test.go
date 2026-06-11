package service

import (
	"testing"

	"db-designer-vkr/internal/model"
)

func TestBuildTransformationsExplainsLogicalAndPhysicalMapping(t *testing.T) {
	steps := BuildTransformations(
		model.Database{Name: "Shop"},
		[]model.Entity{
			{Name: "Customer", Attributes: []model.Attribute{{Name: "email", Type: "VARCHAR(255)", Required: true, Unique: true}}},
			{Name: "Order"},
		},
		[]model.Relation{{From: "Order", To: "Customer", Type: "belongs_to", Cardinality: "many-to-one"}},
	)

	for _, rule := range []string{"database_name_to_create_database", "entity_to_table", "attribute_to_column", "relation_to_foreign_key", "foreign_key_index"} {
		if !hasTransformationRule(steps, rule) {
			t.Fatalf("expected transformation rule %q, got %#v", rule, steps)
		}
	}
}

func hasTransformationRule(steps []model.TransformationStep, rule string) bool {
	for _, step := range steps {
		if step.Rule == rule {
			return true
		}
	}
	return false
}
