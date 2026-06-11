package service

import (
	"encoding/json"
	"os"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestUIExamplesRelationQuality(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := loadUIDomainExamples(t)
	expectedRelations := map[string][]model.Relation{
		"library": {
			{From: "Автор", To: "Книга", Cardinality: "one-to-many"},
			{From: "Категория", To: "Книга", Cardinality: "one-to-many"},
			{From: "Читатель", To: "Выдача", Cardinality: "one-to-many"},
			{From: "Выдача", To: "Книга", Cardinality: "many-to-one"},
			{From: "Бронирование", To: "Книга", Cardinality: "many-to-one"},
			{From: "Выдача", To: "Штраф", Cardinality: "one-to-many"},
		},
		"university": {
			{From: "Факультет", To: "Кафедра", Cardinality: "one-to-many"},
			{From: "Кафедра", To: "Курс", Cardinality: "one-to-many"},
			{From: "Студент", To: "Группа", Cardinality: "many-to-one"},
			{From: "Студент", To: "Курс", Cardinality: "many-to-many"},
			{From: "Преподаватель", To: "Курс", Cardinality: "many-to-many"},
			{From: "Студент", To: "Экзамен", Cardinality: "many-to-many"},
			{From: "Экзамен", To: "Курс", Cardinality: "many-to-one"},
			{From: "Преподаватель", To: "Результат", Cardinality: "one-to-many"},
			{From: "Результат", To: "Экзамен", Cardinality: "many-to-one"},
		},
		"control": {
			{From: "Клиент", To: "Заказ", Cardinality: "one-to-many"},
			{From: "Заказ", To: "Товар", Cardinality: "one-to-many"},
			{From: "Товар", To: "Категория", Cardinality: "many-to-one"},
			{From: "Поставщик", To: "Товар", Cardinality: "one-to-many"},
			{From: "Заказ", To: "Оплата", Cardinality: "one-to-many"},
			{From: "Заказ", To: "Доставка", Cardinality: "one-to-many"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			result, err := AnalyzeTextWithDatabase(scenario.text, scenario.database)
			if err != nil {
				t.Fatalf("AnalyzeText returned error: %v", err)
			}

			for _, entity := range scenario.expectedEntities {
				if !hasEntity(result.Entities, entity) {
					t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
				}
			}
			for _, relation := range expectedRelations[scenario.name] {
				if !hasRelation(result.Relations, relation.From, relation.To, relation.Cardinality) {
					t.Fatalf("expected relation %#v, got %#v", relation, result.Relations)
				}
			}
			if isolated := isolatedEntities(result.Entities, result.Relations); len(isolated) > 0 {
				t.Fatalf("expected no isolated entities, got %v from relations %#v", isolated, result.Relations)
			}
		})
	}
}

type uiDomainExample struct {
	name             string
	text             string
	database         model.Database
	expectedEntities []string
}

func loadUIDomainExamples(t *testing.T) []uiDomainExample {
	t.Helper()

	content, err := os.ReadFile("../../web/domain-examples.json")
	if err != nil {
		t.Fatalf("failed to read UI domain examples: %v", err)
	}

	var raw map[string]struct {
		Text             string   `json:"text"`
		Database         string   `json:"database"`
		ExpectedEntities []string `json:"expectedEntities"`
	}
	if err := json.Unmarshal(content, &raw); err != nil {
		t.Fatalf("failed to parse UI domain examples: %v", err)
	}

	var result []uiDomainExample
	for _, name := range []string{"library", "university", "control"} {
		profile, ok := raw[name]
		if !ok {
			t.Fatalf("missing UI domain example %q", name)
		}
		result = append(result, uiDomainExample{
			name:             name,
			text:             profile.Text,
			database:         model.Database{Name: profile.Database},
			expectedEntities: profile.ExpectedEntities,
		})
	}
	return result
}

func TestFreeFormLibraryDescription(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Мне нужна база для библиотеки. В ней есть читатели, книги, авторы и выдачи. Читатели берут книги. Выдача связана с книгой.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if result.Database.Name != "Библиотека" {
		t.Fatalf("expected library database name, got %#v", result.Database)
	}

	for _, entity := range []string{"Читатель", "Книга", "Автор", "Выдача"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}
	for _, artifact := range []string{"Мне", "Ней", "Нужна", "Берут"} {
		if hasEntity(result.Entities, artifact) {
			t.Fatalf("unexpected artifact entity %q, got %#v", artifact, result.Entities)
		}
	}
	for _, relation := range []model.Relation{
		{From: "Читатель", To: "Книга", Cardinality: "many-to-many"},
		{From: "Выдача", To: "Книга", Cardinality: "many-to-one"},
	} {
		if !hasRelation(result.Relations, relation.From, relation.To, relation.Cardinality) {
			t.Fatalf("expected relation %#v, got %#v", relation, result.Relations)
		}
	}
}

func TestFreeFormDescriptionSkipsRussianConnectors(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	input := "Мне нужна база для библиотеки. " +
		"В ней есть читатели, книги, авторы, а также выдачи. " +
		"Читатель имеет имя, email, а также телефон. " +
		"Книга имеет название, год, а также isbn. " +
		"Еще есть бронирования. " +
		"Бронирование тоже связано с книгой. " +
		"Читатели берут книги."

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{
		"Читатель",
		"Книга",
		"Автор",
		"Выдача",
		"Бронирование",
	} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, entity := range result.Entities {
		for _, artifact := range []string{
			"а",
			"также",
			"еще",
			"ещё",
			"тоже",
			"плюс",
		} {
			if hasAttribute(entity, artifact) {
				t.Fatalf("connector %q must not become an attribute of %q: %#v", artifact, entity.Name, entity.Attributes)
			}
		}
	}

	if !hasRelation(result.Relations, "Читатель", "Книга", "many-to-many") {
		t.Fatalf("expected Reader -> Book relation, got %#v", result.Relations)
	}
	if !hasRelation(result.Relations, "Бронирование", "Книга", "many-to-one") {
		t.Fatalf("expected Reservation -> Book relation, got %#v", result.Relations)
	}
}

func isolatedEntities(entities []model.Entity, relations []model.Relation) []string {
	connected := make(map[string]bool)
	for _, relation := range relations {
		connected[relation.From] = true
		connected[relation.To] = true
	}

	var isolated []string
	for _, entity := range entities {
		if !connected[entity.Name] {
			isolated = append(isolated, entity.Name)
		}
	}
	return isolated
}
