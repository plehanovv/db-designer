package service

import "testing"

func TestValidateSQLPassesGeneratedDDL(t *testing.T) {
	diagnostics := ValidateSQL(`CREATE TABLE customer (
    id SERIAL PRIMARY KEY
);

CREATE TABLE order_table (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER,
    CONSTRAINT fk_order_table_customer FOREIGN KEY (customer_id) REFERENCES customer(id)
);

CREATE INDEX idx_order_table_customer_id ON order_table(customer_id);
`)

	if !containsDiagnostic(diagnostics, "passed internal DDL sanity checks") {
		t.Fatalf("expected success diagnostic, got %#v", diagnostics)
	}
}

func TestValidateSQLDetectsMissingReferencedTable(t *testing.T) {
	diagnostics := ValidateSQL(`CREATE TABLE order_table (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER,
    CONSTRAINT fk_order_table_customer FOREIGN KEY (customer_id) REFERENCES customer(id)
);
`)

	if !containsDiagnostic(diagnostics, "references missing table") {
		t.Fatalf("expected missing table diagnostic, got %#v", diagnostics)
	}
}

func TestValidateSQLDetectsUnbalancedParentheses(t *testing.T) {
	diagnostics := ValidateSQL(`CREATE TABLE customer (
    id SERIAL PRIMARY KEY;
`)

	if !containsDiagnostic(diagnostics, "unbalanced parentheses") {
		t.Fatalf("expected parentheses diagnostic, got %#v", diagnostics)
	}
}
