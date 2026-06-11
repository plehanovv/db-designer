package service

import (
	"strings"
	"testing"

	"db-designer-vkr/internal/model"
)

func TestSemanticAuditHotelBooking(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных гостиницы. Гость имеет имя телефон email. Номер имеет номер этаж цену. Бронирование имеет дату заезда дату выезда статус. Гость оформляет бронирования. Бронирование связано с номером.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Гость", "Номер", "Бронирование"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	if result.Database.Name != "Гостиница" {
		t.Fatalf("expected hotel database name, got %#v", result.Database)
	}

	for _, fragment := range []string{"CREATE DATABASE gostinica", "data_zaezd DATE", "data_vyezd DATE", "gost_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditTicketTracking(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Система учета заявок. Пользователь имеет имя email роль. Сотрудник имеет имя должность. Заявка имеет тему описание статус дату создания. Заявка назначена сотруднику. Сотрудник обрабатывает заявки.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Пользователь", "Сотрудник", "Заявка"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	if result.Database.Name != "Учет Заявка" {
		t.Fatalf("expected ticket tracking database name, got %#v", result.Database)
	}

	for _, fragment := range []string{"CREATE DATABASE uchet_zayavka", "data_sozdanie DATE", "sotrudnik_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditCarRental(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных аренды автомобилей. Клиент имеет имя телефон паспорт. Автомобиль имеет марку модель госномер цену. Договор имеет дату начала дату окончания сумму. Клиент арендует автомобили. Договор связан с автомобилем.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Клиент", "Автомобиль", "Договор"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, fragment := range []string{"CREATE DATABASE arenda_avtomobil", "data_nachalo DATE", "data_okonchanie DATE", "klient_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditClinic(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Предметная область клиники. Пациент имеет имя дату рождения телефон. Врач имеет имя специальность. Прием имеет дату время диагноз. Пациент имеет несколько приемов. Прием назначен врачу.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Пациент", "Врач", "Прием"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, fragment := range []string{"CREATE DATABASE klinika", "data_rozhdenie DATE", "vremya TIME", "pacient_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditInternetShopFreeForm(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных интернет магазина. У нас есть товары которые имеют цену скидку количество на складе. Пользователи у которых есть номер телефона почта дата последней покупки. История заказов которая связана с пользователями.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Товар", "Пользователь", "История"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, fragment := range []string{"CREATE DATABASE internet_magazin", "cena NUMERIC(12,2)", "skidka NUMERIC(12,2)", "kolichestvo INTEGER", "telefon VARCHAR(20)", "pochta VARCHAR(255)", "data_pokupka DATE", "polzovatel_id INTEGER"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditDocumentControlExample(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("Database of educational process. Student has name email. Course has title code. Teacher has name email. Student enrolls in courses. Teacher teaches courses.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Student", "Course", "Teacher"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, relation := range []struct {
		from string
		to   string
	}{
		{"Student", "Course"},
		{"Teacher", "Course"},
	} {
		if !hasManyToManyRelation(result.Relations, relation.from, relation.to) {
			t.Fatalf("expected many-to-many relation %s -> %s, got %#v", relation.from, relation.to, result.Relations)
		}
	}

	for _, fragment := range []string{"CREATE DATABASE educational_process", "CREATE TABLE student_course", "CREATE TABLE teacher_course", "student_id INTEGER NOT NULL", "teacher_id INTEGER NOT NULL", "course_id INTEGER NOT NULL"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func TestSemanticAuditTypicalBusinessRelations(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	result, err := AnalyzeText("База данных компании. Сотрудник имеет имя email. Отдел имеет название. Сотрудник работает в отделе. Пользователь создает заявки. Заявка имеет статус описание. Товар относится к категории. Категория включает товары.")
	if err != nil {
		t.Fatalf("AnalyzeText returned error: %v", err)
	}

	for _, entity := range []string{"Сотрудник", "Отдел", "Пользователь", "Заявка", "Товар", "Категория"} {
		if !hasEntity(result.Entities, entity) {
			t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
		}
	}

	for _, expected := range []struct {
		from        string
		to          string
		cardinality string
	}{
		{"Сотрудник", "Отдел", "many-to-one"},
		{"Пользователь", "Заявка", "one-to-many"},
		{"Товар", "Категория", "many-to-one"},
		{"Категория", "Товар", "one-to-many"},
	} {
		if !hasRelation(result.Relations, expected.from, expected.to, expected.cardinality) {
			t.Fatalf("expected relation %#v, got %#v", expected, result.Relations)
		}
	}

	for _, fragment := range []string{"otdel_id INTEGER", "polzovatel_id INTEGER", "kategoriya_id INTEGER", "CREATE INDEX idx_sotrudnik_otdel_id"} {
		if !strings.Contains(result.SQL, fragment) {
			t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
		}
	}
}

func hasManyToManyRelation(relations []model.Relation, from string, to string) bool {
	for _, relation := range relations {
		if relation.From == from && relation.To == to && relation.Cardinality == "many-to-many" {
			return true
		}
	}
	return false
}

func hasRelation(relations []model.Relation, from string, to string, cardinality string) bool {
	for _, relation := range relations {
		if relation.From == from && relation.To == to && relation.Cardinality == cardinality {
			return true
		}
	}
	return false
}
