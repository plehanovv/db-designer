package service

import (
	"strings"
	"testing"
)

func TestAnalyzeTextRussianDomainDescription(t *testing.T) {
	input := "\u041a\u043b\u0438\u0435\u043d\u0442 \u0438\u043c\u0435\u0435\u0442 \u0438\u043c\u044f email \u0442\u0435\u043b\u0435\u0444\u043e\u043d. " +
		"\u0417\u0430\u043a\u0430\u0437 \u0438\u043c\u0435\u0435\u0442 \u043d\u043e\u043c\u0435\u0440 \u0434\u0430\u0442\u0443 \u0441\u0443\u043c\u043c\u0443. " +
		"\u0417\u0430\u043a\u0430\u0437 \u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0438\u0442 \u043a\u043b\u0438\u0435\u043d\u0442\u0443."

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Entities) != 2 {
		t.Fatalf("expected 2 entities, got %d: %#v", len(result.Entities), result.Entities)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d: %#v", len(result.Relations), result.Relations)
	}

	if result.Relations[0].From != "\u0417\u0430\u043a\u0430\u0437" || result.Relations[0].To != "\u041a\u043b\u0438\u0435\u043d\u0442" {
		t.Fatalf("unexpected relation: %#v", result.Relations[0])
	}

	for _, fragment := range []string{
		"CREATE TABLE zakaz",
		"klient_id INTEGER",
		"FOREIGN KEY (klient_id) REFERENCES klient(id)",
		"summa NUMERIC(12,2)",
	} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}
