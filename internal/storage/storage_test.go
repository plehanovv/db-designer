package storage

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestBuildSaveSQLIncludesProjectArtifacts(t *testing.T) {
	sql, err := buildSaveSQL("abc123", AnalysisRecord{
		InputText: "Customer's order",
		InputType: "text",
		Database:  model.Database{Name: "Shop"},
		Entities: []model.Entity{
			{
				Name: "Customer",
				Attributes: []model.Attribute{
					{Name: "email", Type: "VARCHAR(255)", Required: true, Unique: true},
				},
			},
			{Name: "Order"},
		},
		Relations: []model.Relation{
			{From: "Order", To: "Customer", Type: "belongs_to", Cardinality: "many-to-one"},
		},
		SQL: "CREATE TABLE customer (id SERIAL PRIMARY KEY);",
	})
	if err != nil {
		t.Fatalf("buildSaveSQL returned error: %v", err)
	}

	for _, fragment := range []string{
		"INSERT INTO text_description",
		"Customer''s order",
		"INSERT INTO entity",
		"INSERT INTO attribute",
		"INSERT INTO relation",
		"INSERT INTO database_schema",
		"belongs_to",
		"many-to-one",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, sql)
		}
	}
}

func TestNewFromEnvDisabledWithoutDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_DESIGNER_DATABASE_URL", "")

	store, enabled := NewFromEnv()
	if enabled || store != nil {
		t.Fatalf("expected storage to be disabled without database URL")
	}
}
