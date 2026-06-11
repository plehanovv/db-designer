package analyzer

import "testing"

func TestExtractDatabaseFromRussianIntro(t *testing.T) {
	database := ExtractDatabase("База данных интернет магазина. У нас есть товары.")

	if database.Name != "Интернет Магазин" {
		t.Fatalf("expected Internet shop database name, got %#v", database)
	}
}

func TestExtractDatabaseFallback(t *testing.T) {
	database := ExtractDatabase("Клиент имеет имя.")

	if database.Name != "database" {
		t.Fatalf("expected default database name, got %#v", database)
	}
}
