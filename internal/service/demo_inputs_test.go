package service

import (
	"os"
	"strings"
	"testing"
)

func TestDemoInputFiles(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := []struct {
		name         string
		path         string
		minEntities  int
		minRelations int
	}{
		{name: "text", path: "../../examples/domain_description.txt", minEntities: 3, minRelations: 2},
		{name: "university_text", path: "../../examples/university_process.txt", minEntities: 8, minRelations: 8},
		{name: "json", path: "../../examples/structured_model.json", minEntities: 3, minRelations: 2},
		{name: "csv", path: "../../examples/structured_model.csv", minEntities: 3, minRelations: 2},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			content, err := os.ReadFile(scenario.path)
			if err != nil {
				t.Fatalf("failed to read demo file: %v", err)
			}

			result, err := AnalyzeText(string(content))
			if err != nil {
				t.Fatalf("AnalyzeText returned error: %v", err)
			}

			if len(result.Entities) < scenario.minEntities {
				t.Fatalf("expected at least %d entities, got %d: %#v", scenario.minEntities, len(result.Entities), result.Entities)
			}
			if len(result.Relations) < scenario.minRelations {
				t.Fatalf("expected at least %d relations, got %d: %#v", scenario.minRelations, len(result.Relations), result.Relations)
			}
			if result.SQL == "" {
				t.Fatal("expected generated SQL")
			}
			if scenario.name == "text" {
				if hasEntity(result.Entities, "Isbn") {
					t.Fatalf("ISBN must remain a scalar attribute, got entities: %#v", result.Entities)
				}
				book := findEntity(result.Entities, "Книга")
				if !hasAttribute(book, "isbn") {
					t.Fatalf("expected book isbn attribute, got: %#v", book.Attributes)
				}
			}
		})
	}
}

func TestLogisticsDemoInputProducesUsableSchema(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	input := "Предметная область описывает работу распределительного логистического центра. " +
		"В системе есть клиенты, заказы, грузы, склады, зоны склада, ячейки хранения, водители, транспортные средства, маршруты, рейсы, накладные, платежи, страховые полисы, инциденты и сотрудники. " +
		"У клиента есть название, ИНН, телефон, email и адрес. " +
		"Клиент оформляет заказы. " +
		"У заказа есть номер, дата создания, статус и стоимость доставки. " +
		"Заказ содержит несколько грузов. " +
		"У груза есть код, вес, объем и температурный режим. " +
		"Груз хранится на складе. " +
		"Склад содержит зоны склада. " +
		"Зона содержит ячейки хранения. " +
		"Груз размещается в ячейке склада. " +
		"Заказ связан с накладной. " +
		"Заказ связан с платежом. " +
		"Груз связан со страховым полисом. " +
		"Водитель выполняет рейсы. " +
		"Рейс использует транспортное средство. " +
		"Рейс выполняется по маршруту. " +
		"Маршрут содержит пункты маршрута. " +
		"Во время рейса возникают инциденты. " +
		"Сотрудник обрабатывает заказы."

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{
		"Клиент", "Заказ", "Груз", "Склад", "Зона", "Ячейка", "Водитель",
		"Транспортное средство", "Маршрут", "Рейс", "Накладная", "Платеж",
		"Полис", "Инцидент", "Сотрудник",
	} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got: %#v", entity, result.Entities)
		}
	}

	for _, noise := range []string{
		"Инн", "Код", "Площадь", "Режим", "Вместимость", "Центр",
		"Страхов", "Транспортн",
	} {
		if hasEntity(result.Entities, noise) {
			t.Fatalf("scalar attribute/noise %q must not become entity: %#v", noise, result.Entities)
		}
	}

	for _, relation := range []struct {
		from        string
		to          string
		cardinality string
	}{
		{from: "Клиент", to: "Заказ", cardinality: "one-to-many"},
		{from: "Заказ", to: "Груз", cardinality: "one-to-many"},
		{from: "Груз", to: "Склад", cardinality: "many-to-one"},
		{from: "Склад", to: "Зона", cardinality: "one-to-many"},
		{from: "Зона", to: "Ячейка", cardinality: "one-to-many"},
		{from: "Груз", to: "Ячейка", cardinality: "many-to-one"},
		{from: "Водитель", to: "Рейс", cardinality: "one-to-many"},
		{from: "Рейс", to: "Транспортное средство", cardinality: "many-to-one"},
		{from: "Рейс", to: "Маршрут", cardinality: "many-to-one"},
		{from: "Рейс", to: "Инцидент", cardinality: "one-to-many"},
		{from: "Сотрудник", to: "Заказ", cardinality: "one-to-many"},
	} {
		if !hasRelation(result.Relations, relation.from, relation.to, relation.cardinality) {
			t.Fatalf("expected relation %s -> %s (%s), got: %#v", relation.from, relation.to, relation.cardinality, result.Relations)
		}
	}

	for _, table := range []string{
		"CREATE TABLE klient",
		"CREATE TABLE zakaz",
		"CREATE TABLE gruz",
		"CREATE TABLE sklad",
		"CREATE TABLE reys",
	} {
		if !strings.Contains(result.SQL, table) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", table, result.SQL)
		}
	}
}

func TestSemiFormalLogisticsInputKeepsDriverFieldsAsAttributes(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	input := "Предметная область: складская логистика.\n\n" +
		"У клиента есть название, телефон, email и адрес.\n" +
		"У заказа есть номер, дата, статус и стоимость.\n" +
		"У груза есть код, название, вес и объем.\n" +
		"У склада есть название, адрес и площадь.\n" +
		"У ячейки есть код, вместимость и статус.\n" +
		"У водителя есть ФИО, телефон и категория.\n" +
		"У рейса есть номер, дата, статус и расстояние.\n" +
		"У маршрута есть название, начальный пункт и конечный пункт.\n" +
		"У платежа есть номер, дата, сумма и статус.\n\n" +
		"Клиент оформляет заказы.\n" +
		"Заказ содержит несколько грузов.\n" +
		"Груз хранится на складе.\n" +
		"Склад содержит несколько ячеек.\n" +
		"Груз размещается в ячейке.\n" +
		"Водитель выполняет рейсы.\n" +
		"Рейс выполняется по маршруту.\n" +
		"Заказ связан с платежом."

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if !hasEntity(result.Entities, "Водитель") {
		t.Fatalf("expected driver entity, got: %#v", result.Entities)
	}
	for _, noise := range []string{"Фио", "Категория"} {
		if hasEntity(result.Entities, noise) {
			t.Fatalf("%q must remain a driver attribute, got entities: %#v", noise, result.Entities)
		}
	}

	driver := findEntity(result.Entities, "Водитель")
	for _, attribute := range []string{"фио", "телефон", "категория"} {
		if !hasAttribute(driver, attribute) {
			t.Fatalf("expected driver attribute %q, got: %#v", attribute, driver.Attributes)
		}
	}

	for _, relation := range []struct {
		from        string
		to          string
		cardinality string
	}{
		{from: "Клиент", to: "Заказ", cardinality: "one-to-many"},
		{from: "Заказ", to: "Груз", cardinality: "one-to-many"},
		{from: "Груз", to: "Склад", cardinality: "many-to-one"},
		{from: "Склад", to: "Ячейка", cardinality: "one-to-many"},
		{from: "Груз", to: "Ячейка", cardinality: "many-to-one"},
		{from: "Водитель", to: "Рейс", cardinality: "one-to-many"},
		{from: "Рейс", to: "Маршрут", cardinality: "many-to-one"},
		{from: "Заказ", to: "Платеж", cardinality: "many-to-one"},
	} {
		if !hasRelation(result.Relations, relation.from, relation.to, relation.cardinality) {
			t.Fatalf("expected relation %s -> %s (%s), got: %#v", relation.from, relation.to, relation.cardinality, result.Relations)
		}
	}
}
