package service

import (
	"strings"
	"testing"
)

func TestExtendedDomainSemantics(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := []struct {
		name      string
		input     string
		entities  []string
		fragments []string
	}{
		{
			name:      "insurance",
			input:     "Система страхования. Клиент имеет имя телефон паспорт. Агент имеет имя телефон. Полис имеет номер дату премию статус. Случай имеет дату описание сумму. Выплата имеет сумму дату статус. Клиент оформляет полисы. Агент страхует клиентов. Клиент подает случаи. Выплата связана со случаем.",
			entities:  []string{"Клиент", "Агент", "Полис", "Случай", "Выплата"},
			fragments: []string{"CREATE DATABASE strahovanie", "premiya NUMERIC(12,2)", "summa NUMERIC(12,2)", "klient_id INTEGER", "sluchay_id INTEGER"},
		},
		{
			name:      "airline",
			input:     "Система авиаперевозок. Пассажир имеет имя паспорт email. Рейс имеет номер дату время маршрут статус. Билет имеет номер место цену. Аэропорт имеет название город. Пассажир покупает билеты. Билет связан с рейсом. Рейс связан с аэропортом.",
			entities:  []string{"Пассажир", "Рейс", "Билет", "Аэропорт"},
			fragments: []string{"CREATE DATABASE aviaperevozka", "marshrut VARCHAR(255)", "mesto INTEGER", "passazhir_id INTEGER", "reys_id INTEGER", "aeroport_id INTEGER"},
		},
		{
			name:      "fitness",
			input:     "База данных фитнеса. Клиент имеет имя телефон email. Абонемент имеет название дату начала дату окончания цену статус. Тренер имеет имя телефон специальность. Зал имеет номер вместимость. Занятие имеет дату время название. Клиент покупает абонементы. Клиент посещает занятия. Занятие связано с тренером. Занятие связано с залом.",
			entities:  []string{"Клиент", "Абонемент", "Тренер", "Зал", "Занятие"},
			fragments: []string{"CREATE DATABASE fitnes", "data_nachalo DATE", "data_okonchanie DATE", "cena NUMERIC(12,2)", "klient_id INTEGER", "trener_id INTEGER", "zal_id INTEGER"},
		},
		{
			name:      "exams",
			input:     "Система экзаменов. Ученик имеет имя класс email. Экзамен имеет название дату кабинет. Преподаватель имеет имя email. Результат имеет оценку дату статус. Ученик сдает экзамены. Преподаватель оценивает результаты. Результат связан с учеником. Результат связан с экзаменом.",
			entities:  []string{"Ученик", "Экзамен", "Преподаватель", "Результат"},
			fragments: []string{"CREATE DATABASE ekzamen", "klass VARCHAR(255)", "kabinet VARCHAR(255)", "ocenka INTEGER", "uchenik_id INTEGER", "ekzamen_id INTEGER"},
		},
		{
			name:      "document-workflow",
			input:     "Система документооборота. Документ имеет номер дату тему статус. Сотрудник имеет имя должность email. Согласование имеет дату статус комментарий. Подпись имеет дату статус. Сотрудник создает документы. Сотрудник согласует документы. Согласование связано с документом. Подпись связана с документом.",
			entities:  []string{"Документ", "Сотрудник", "Согласование", "Подпись"},
			fragments: []string{"CREATE DATABASE dokumentooborot", "tema VARCHAR(255)", "kommentariy VARCHAR(255)", "sotrudnik_id INTEGER", "dokument_id INTEGER"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			result, err := AnalyzeText(scenario.input)
			if err != nil {
				t.Fatalf("AnalyzeText returned error: %v", err)
			}

			for _, entity := range scenario.entities {
				if !hasEntity(result.Entities, entity) {
					t.Fatalf("expected entity %q, got %#v", entity, result.Entities)
				}
			}

			for _, fragment := range scenario.fragments {
				if !strings.Contains(result.SQL, fragment) {
					t.Fatalf("expected SQL to contain %q, got:\n%s", fragment, result.SQL)
				}
			}
		})
	}
}
