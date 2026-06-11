package generator

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestGenerateSQLPlacesOneToManyForeignKeyOnTargetTable(t *testing.T) {
	sql := GenerateSQL(
		model.Database{},
		[]model.Entity{
			{Name: "Customer"},
			{Name: "Order"},
		},
		[]model.Relation{
			{From: "Customer", To: "Order", Type: "has", Cardinality: "one-to-many"},
		},
	)

	if strings.Contains(sql, "order_table_id INTEGER") {
		t.Fatalf("one-to-many relation must not place FK on source table, got:\n%s", sql)
	}

	if !strings.Contains(sql, "customer_id INTEGER") ||
		!strings.Contains(sql, "FOREIGN KEY (customer_id) REFERENCES customer(id)") {
		t.Fatalf("one-to-many relation must place FK on target table, got:\n%s", sql)
	}
}

func TestGenerateSQLCreatesJunctionTableForManyToMany(t *testing.T) {
	sql := GenerateSQL(
		model.Database{},
		[]model.Entity{
			{Name: "Student"},
			{Name: "Course"},
		},
		[]model.Relation{
			{From: "Student", To: "Course", Type: "associated_with", Cardinality: "many-to-many"},
		},
	)

	for _, fragment := range []string{
		"CREATE TABLE student_course",
		"student_id INTEGER NOT NULL",
		"course_id INTEGER NOT NULL",
		"PRIMARY KEY (student_id, course_id)",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, sql)
		}
	}
}

func TestGenerateSQLCreatesDatabaseWhenDomainIsKnown(t *testing.T) {
	sql := GenerateSQL(
		model.Database{Name: "Интернет магазин"},
		[]model.Entity{{Name: "Товар"}},
		nil,
	)

	if !strings.Contains(sql, "CREATE DATABASE internet_magazin") {
		t.Fatalf("expected database creation statement, got:\n%s", sql)
	}
}

func TestGenerateSQLCreatesAttributeConstraintsAndForeignKeyIndexes(t *testing.T) {
	sql := GenerateSQL(
		model.Database{},
		[]model.Entity{
			{
				Name: "Customer",
				Attributes: []model.Attribute{
					{Name: "email", Type: "VARCHAR(255)", Required: true, Unique: true},
				},
			},
			{Name: "Order"},
		},
		[]model.Relation{
			{From: "Order", To: "Customer", Type: "belongs_to", Cardinality: "many-to-one"},
		},
	)

	for _, fragment := range []string{
		"email VARCHAR(255) NOT NULL UNIQUE",
		"customer_id INTEGER",
		"CREATE INDEX idx_order_table_customer_id ON order_table(customer_id);",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, sql)
		}
	}
}
