package service

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestValidateModelDetectsPracticalSchemaProblems(t *testing.T) {
	diagnostics := ValidateModel(
		model.Database{Name: ""},
		[]model.Entity{
			{
				Name: "Order",
				Attributes: []model.Attribute{
					{Name: "email", Type: "VARCHAR(255)"},
					{Name: "email", Type: "VARCHAR(255)"},
					{Name: "custom", Type: "UNKNOWN"},
				},
			},
			{Name: "Order", Attributes: []model.Attribute{}},
		},
		[]model.Relation{
			{From: "Order", To: "Customer", Type: "belongs_to", Cardinality: "many-to-one"},
			{From: "Order", To: "Order", Type: "associated_with", Cardinality: "strange"},
		},
	)

	for _, expected := range []string{
		"Database name is empty",
		"Duplicate entity name",
		"duplicate attribute",
		"unsupported type",
		"does not match any entity",
		"unsupported cardinality",
	} {
		if !containsDiagnostic(diagnostics, expected) {
			t.Fatalf("expected diagnostic containing %q, got %#v", expected, diagnostics)
		}
	}
}

func TestValidateModelWarnsAboutIncompleteCurrentResult(t *testing.T) {
	diagnostics := ValidateModel(
		model.Database{Name: "Library"},
		[]model.Entity{
			{Name: "Reader"},
			{Name: "Book"},
		},
		nil,
	)

	for _, expected := range []string{
		`Entity "Reader" has no scalar attributes`,
		`Entity "Book" has no scalar attributes`,
		"Several entities were detected, but no relations were found",
	} {
		if !containsDiagnostic(diagnostics, expected) {
			t.Fatalf("expected diagnostic containing %q, got %#v", expected, diagnostics)
		}
	}
}

func containsDiagnostic(diagnostics []model.Diagnostic, fragment string) bool {
	for _, diagnostic := range diagnostics {
		if strings.Contains(diagnostic.Message, fragment) {
			return true
		}
	}
	return false
}
