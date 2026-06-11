package service

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestAnalyzeTextStructuredJSONInput(t *testing.T) {
	input := `{
		"database": {"name": "Shop"},
		"entities": [
			{"name": "Customer", "attributes": [{"name": "email", "type": "VARCHAR(255)", "required": true}]},
			{"name": "Order", "attributes": [{"name": "amount", "type": "NUMERIC(12,2)", "required": false}]}
		],
		"relations": [
			{"from": "Order", "to": "Customer", "type": "belongs_to", "cardinality": "many-to-one"}
		]
	}`

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if result.Database.Name != "Shop" || len(result.Entities) != 2 || len(result.Relations) != 1 {
		t.Fatalf("unexpected structured model: %#v", result)
	}

	for _, fragment := range []string{"CREATE DATABASE shop", "CREATE TABLE order_table", "customer_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}

	if !containsDiagnostic(result.Diagnostics, "Structured JSON input") {
		t.Fatalf("expected structured input diagnostic, got %#v", result.Diagnostics)
	}
}

func TestAnalyzeTextInvalidStructuredJSONInput(t *testing.T) {
	_, err := AnalyzeText(`{"entities": [}`)
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestAnalyzeTextStructuredJSONGeneratesTransformations(t *testing.T) {
	input := `{
		"database": {"name": "Library"},
		"entities": [
			{"name": "Reader", "attributes": [{"name": "email", "type": "VARCHAR(255)", "required": true}]},
			{"name": "Book", "attributes": [{"name": "title", "type": "TEXT"}]}
		],
		"relations": [
			{"from": "Reader", "to": "Book", "type": "associated_with", "cardinality": "many-to-many"}
		]
	}`

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Transformations) == 0 {
		t.Fatal("expected transformation steps for structured JSON model")
	}
	if !strings.Contains(result.SQL, "CREATE TABLE reader_book") {
		t.Fatalf("expected many-to-many junction table in SQL, got:\n%s", result.SQL)
	}
}

func TestAnalyzeTextStructuredCSVInput(t *testing.T) {
	input := `kind,database,entity,attribute,type,required,unique,from,to,relation_type,cardinality
attribute,Shop,Customer,email,VARCHAR(255),true,true,,,, 
attribute,Shop,Order,amount,"NUMERIC(12,2)",false,false,,,, 
relation,Shop,,,,,,Order,Customer,belongs_to,many-to-one`

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if result.Database.Name != "Shop" || len(result.Entities) != 2 || len(result.Relations) != 1 {
		t.Fatalf("unexpected CSV model: %#v", result)
	}

	for _, fragment := range []string{"CREATE DATABASE shop", "email VARCHAR(255) NOT NULL UNIQUE", "customer_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}

	if !containsDiagnostic(result.Diagnostics, "Structured CSV input") {
		t.Fatalf("expected CSV diagnostic, got %#v", result.Diagnostics)
	}
}

func TestAnalyzeTextStructuredCSVDatabaseCanAppearAfterFirstRow(t *testing.T) {
	input := `kind,database,entity,attribute,type,required,unique,from,to,relation_type,cardinality
relation,,,,,,,Order,Customer,belongs_to,many-to-one
attribute,Shop,Customer,email,VARCHAR(255),true,true,,,, 
attribute,,Order,amount,"NUMERIC(12,2)",false,false,,,,`

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if result.Database.Name != "Shop" {
		t.Fatalf("expected database name from later CSV row, got %q", result.Database.Name)
	}
	if !strings.Contains(result.SQL, "CREATE DATABASE shop") {
		t.Fatalf("expected SQL to contain database creation, got:\n%s", result.SQL)
	}
	if !hasCandidateRule(result.Explanation.Candidates, "structured_csv_entity") {
		t.Fatalf("expected CSV-specific candidate rule, got %#v", result.Explanation.Candidates)
	}
}

func hasCandidateRule(candidates []model.Candidate, rule string) bool {
	for _, candidate := range candidates {
		if candidate.Rule == rule {
			return true
		}
	}
	return false
}
