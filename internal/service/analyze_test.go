package service

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestAnalyzeTextRussianDomainDescription(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	input := "Клиент имеет имя email телефон. " +
		"Заказ имеет номер дату сумму. " +
		"Заказ принадлежит клиенту."

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

	if result.Relations[0].From != "Заказ" || result.Relations[0].To != "Клиент" {
		t.Fatalf("unexpected relation: %#v", result.Relations[0])
	}

	if len(result.Explanation.Candidates) == 0 {
		t.Fatal("expected explanation candidates")
	}

	if !hasCandidate(result.Explanation.Candidates, "relation", "Заказ", "Клиент") {
		t.Fatalf("expected relation candidate in explanation, got: %#v", result.Explanation.Candidates)
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

func TestAnalyzeTextDoesNotPromoteAttributeToEntity(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Customer has apple email phone. Order has number date amount. Order belongs to customer.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Entities) != 2 {
		t.Fatalf("expected only Customer and Order entities, got %d: %#v", len(result.Entities), result.Entities)
	}

	if hasEntity(result.Entities, "Apple") {
		t.Fatalf("attribute apple must not be promoted to entity: %#v", result.Entities)
	}

	if strings.Contains(result.SQL, "CREATE TABLE apple") || strings.Contains(result.SQL, "apple_id") {
		t.Fatalf("attribute apple must not generate table or foreign key, got:\n%s", result.SQL)
	}

	if !strings.Contains(result.SQL, "apple TEXT") {
		t.Fatalf("expected apple to remain Customer attribute, got:\n%s", result.SQL)
	}
}

func TestAnalyzeTextSkipsAttributeListConnectors(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Customer has an email and phone. Order belongs to customer.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if strings.Contains(result.SQL, " an ") || strings.Contains(result.SQL, " and ") {
		t.Fatalf("attribute connectors must not be emitted as columns, got:\n%s", result.SQL)
	}

	for _, fragment := range []string{"email VARCHAR(255)", "phone VARCHAR(20)", "customer_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestAnalyzeTextRussianPossessiveAttributes(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("У клиента есть имя и телефон.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Entities) != 1 || result.Entities[0].Name != "Клиент" {
		t.Fatalf("expected only Client entity, got: %#v", result.Entities)
	}

	for _, fragment := range []string{"imya VARCHAR(255)", "telefon VARCHAR(20)"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestAnalyzeTextQuantifiedPossessionRelation(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Один клиент может иметь несколько заказов.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected one quantified relation, got: %#v", result.Relations)
	}

	relation := result.Relations[0]
	if relation.From != "Клиент" || relation.To != "Заказ" || relation.Cardinality != "one-to-many" {
		t.Fatalf("unexpected quantified relation: %#v", relation)
	}

	if !strings.Contains(result.SQL, "klient_id INTEGER") || !strings.Contains(result.SQL, "REFERENCES klient(id)") {
		t.Fatalf("one-to-many relation must place FK on target table, got:\n%s", result.SQL)
	}
}

func TestAnalyzeTextEnglishNaturalOneToManyRelation(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Each customer can place many orders.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected one natural language relation, got: %#v", result.Relations)
	}

	relation := result.Relations[0]
	if relation.From != "Customer" || relation.To != "Order" || relation.Cardinality != "one-to-many" {
		t.Fatalf("unexpected natural language relation: %#v", relation)
	}
}

func TestAnalyzeTextEnglishCanHaveManyRelation(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Each customer can have many orders.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected one can-have-many relation, got: %#v", result.Relations)
	}

	relation := result.Relations[0]
	if relation.From != "Customer" || relation.To != "Order" || relation.Cardinality != "one-to-many" {
		t.Fatalf("unexpected can-have-many relation: %#v", relation)
	}
}

func TestAnalyzeTextRussianAssociatedRelation(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Заказ связан с клиентом.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected one associated relation, got: %#v", result.Relations)
	}

	relation := result.Relations[0]
	if relation.From != "Заказ" || relation.To != "Клиент" || relation.Type != "associated_with" {
		t.Fatalf("unexpected associated relation: %#v", relation)
	}
}

func TestAnalyzeTextReturnsDiagnostics(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Customer has email.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Diagnostics) == 0 {
		t.Fatal("expected diagnostics")
	}
}

func TestAnalyzeTextDetectsAttributeConstraints(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Customer has required unique email phone.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	customer := findEntity(result.Entities, "Customer")
	email := findAttribute(customer, "email")
	if email.Name == "" {
		t.Fatalf("expected email attribute, got %#v", customer.Attributes)
	}
	if !email.Required || !email.Unique {
		t.Fatalf("expected email to be required and unique, got %#v", email)
	}
	if hasAttribute(customer, "required") || hasAttribute(customer, "unique") {
		t.Fatalf("constraint markers must not become attributes, got %#v", customer.Attributes)
	}
	if !strings.Contains(result.SQL, "email VARCHAR(255) NOT NULL UNIQUE") {
		t.Fatalf("expected SQL constraint for email, got:\n%s", result.SQL)
	}
}

func TestAnalyzeTextStorageEnumeration(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("В системе хранятся клиенты, заказы и товары.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entityName := range []string{
		"Клиент",
		"Заказ",
		"Товар",
	} {
		if !hasEntity(result.Entities, entityName) {
			t.Fatalf("expected entity %q, got: %#v", entityName, result.Entities)
		}
	}
}

func TestAnalyzeTextRussianActionRelation(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Клиенты оформляют заказы.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if len(result.Relations) != 1 {
		t.Fatalf("expected one action relation, got: %#v", result.Relations)
	}

	relation := result.Relations[0]
	if relation.From != "Клиент" || relation.To != "Заказ" || relation.Cardinality != "one-to-many" {
		t.Fatalf("unexpected action relation: %#v", relation)
	}
}

func TestAnalyzeTextContainmentMixedEntityAndAttributes(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Заказ содержит товары, дату и сумму.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	if !hasEntity(result.Entities, "Заказ") || !hasEntity(result.Entities, "Товар") {
		t.Fatalf("expected Order and Product entities, got: %#v", result.Entities)
	}

	if len(result.Relations) != 1 || result.Relations[0].From != "Заказ" || result.Relations[0].To != "Товар" {
		t.Fatalf("expected Order contains Product relation, got: %#v", result.Relations)
	}

	for _, fragment := range []string{"data DATE", "summa NUMERIC(12,2)"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestAnalyzeTextLibrarySampleDoesNotPromoteCompoundDateParts(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных библиотеки. " +
		"Читатель имеет имя email телефон. " +
		"Книга имеет название автора год. " +
		"Выдача имеет дату выдачи и дату возврата. " +
		"Читатель имеет несколько выдач. " +
		"Выдача связана с книгой.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entityName := range []string{
		"Читатель",
		"Книга",
		"Выдача",
	} {
		if !hasEntity(result.Entities, entityName) {
			t.Fatalf("expected entity %q, got: %#v", entityName, result.Entities)
		}
	}

	for _, entityName := range []string{
		"Библиотек",
		"Возврата",
		"Выдач",
	} {
		if hasEntity(result.Entities, entityName) {
			t.Fatalf("compound date/domain noise %q must not become entity: %#v", entityName, result.Entities)
		}
	}

	loan := findEntity(result.Entities, "Выдача")
	for _, attributeName := range []string{
		"дата_выдача",
		"дата_возврат",
	} {
		if !hasAttribute(loan, attributeName) {
			t.Fatalf("expected loan attribute %q, got: %#v", attributeName, loan.Attributes)
		}
	}

	if !hasRelation(result.Relations, "Выдача", "Книга", "many-to-one") {
		t.Fatalf("expected Loan associated with Book relation, got: %#v", result.Relations)
	}
}

func TestAnalyzeTextHotelSampleDoesNotPromoteCompoundDateParts(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных гостиницы. " +
		"Гость имеет имя телефон email. " +
		"Номер имеет номер этаж цену. " +
		"Бронирование имеет дату заезда дату выезда статус. " +
		"Гость оформляет бронирования. " +
		"Бронирование связано с номером.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entityName := range []string{
		"Заезда",
		"Выезда",
	} {
		if hasEntity(result.Entities, entityName) {
			t.Fatalf("compound date part %q must not become entity: %#v", entityName, result.Entities)
		}
	}
}

func TestAnalyzeTextSemiFreeInternetShopDescription(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	input := "База данных интернет магазина. У нас есть товары которые имеют цену скидку количество на складе. " +
		"Пользователи у которых есть номер телефона почта дата последней покупки. " +
		"История заказов которая связана с пользователями."

	result, err := AnalyzeText(input)
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entityName := range []string{"Товар", "Пользователь", "История"} {
		if !hasEntity(result.Entities, entityName) {
			t.Fatalf("expected entity %q, got: %#v", entityName, result.Entities)
		}
	}

	for _, entityName := range []string{"База", "Интернет", "Магазин", "Который"} {
		if hasEntity(result.Entities, entityName) {
			t.Fatalf("introductory/noise word %q must not become entity: %#v", entityName, result.Entities)
		}
	}

	product := findEntity(result.Entities, "Товар")
	for _, attributeName := range []string{"цена", "скидка", "количество"} {
		if !hasAttribute(product, attributeName) {
			t.Fatalf("expected product attribute %q, got: %#v", attributeName, product.Attributes)
		}
	}

	user := findEntity(result.Entities, "Пользователь")
	for _, attributeName := range []string{"телефон", "почта", "дата_покупка"} {
		if !hasAttribute(user, attributeName) {
			t.Fatalf("expected user attribute %q, got: %#v", attributeName, user.Attributes)
		}
	}

	if len(result.Relations) != 1 || result.Relations[0].From != "История" || result.Relations[0].To != "Пользователь" {
		t.Fatalf("expected History associated with User relation, got: %#v", result.Relations)
	}

	for _, fragment := range []string{
		"CREATE TABLE tovar",
		"cena NUMERIC(12,2)",
		"skidka NUMERIC(12,2)",
		"kolichestvo INTEGER",
		"telefon VARCHAR(20)",
		"pochta VARCHAR(255) NOT NULL",
		"data_pokupka DATE",
	} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func hasEntity(entities []model.Entity, name string) bool {
	for _, entity := range entities {
		if entity.Name == name {
			return true
		}
	}

	return false
}

func findEntity(entities []model.Entity, name string) model.Entity {
	for _, entity := range entities {
		if entity.Name == name {
			return entity
		}
	}

	return model.Entity{}
}

func hasAttribute(entity model.Entity, name string) bool {
	for _, attribute := range entity.Attributes {
		if attribute.Name == name {
			return true
		}
	}

	return false
}

func findAttribute(entity model.Entity, name string) model.Attribute {
	for _, attribute := range entity.Attributes {
		if attribute.Name == name {
			return attribute
		}
	}
	return model.Attribute{}
}

func hasCandidate(candidates []model.Candidate, kind string, owner string, target string) bool {
	for _, candidate := range candidates {
		if candidate.Kind == kind && candidate.Owner == owner && candidate.Target == target {
			return true
		}
	}

	return false
}
