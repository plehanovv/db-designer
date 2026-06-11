package service

import (
	"strings"
	"testing"
)

func TestOntologyExpansionDomains(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := []struct {
		name      string
		input     string
		entities  []string
		fragments []string
	}{
		{
			name:      "pharmacy",
			input:     "База данных аптеки. Клиент имеет имя телефон. Лекарство имеет название дозировку цену остаток. Поставщик имеет имя телефон. Продажа имеет дату сумму количество. Клиент покупает лекарства. Поставщик поставляет лекарства. Продажа связана с клиентом.",
			entities:  []string{"Клиент", "Лекарство", "Поставщик", "Продажа"},
			fragments: []string{"CREATE DATABASE apteka", "dozirovka VARCHAR(255)", "ostatok INTEGER", "klient_id INTEGER", "postavschik_id INTEGER"},
		},
		{
			name:      "cinema",
			input:     "База данных кинотеатра. Фильм имеет название жанр длительность. Зал имеет номер вместимость. Сеанс имеет дату время цену. Билет имеет номер место цену. Сеанс связан с фильмом. Сеанс связан с залом. Билет связан с сеансом.",
			entities:  []string{"Фильм", "Зал", "Сеанс", "Билет"},
			fragments: []string{"CREATE DATABASE kinoteatr", "zhanr VARCHAR(255)", "dlitelnost INTEGER", "film_id INTEGER", "zal_id INTEGER", "seans_id INTEGER"},
		},
		{
			name:      "museum",
			input:     "Система музея. Экспонат имеет название автор год описание. Выставка имеет название дату тему. Зал имеет номер название. Посетитель имеет имя email. Выставка содержит экспонаты. Выставка связана с залом. Посетитель посещает выставки.",
			entities:  []string{"Экспонат", "Выставка", "Зал", "Посетитель"},
			fragments: []string{"CREATE DATABASE muzey", "avtor VARCHAR(255)", "god INTEGER", "tema VARCHAR(255)", "vystavka_id INTEGER", "zal_id INTEGER"},
		},
		{
			name:      "construction",
			input:     "Система строительства. Объект имеет название адрес бюджет статус. Этап имеет название дату начала дату окончания стоимость. Подрядчик имеет имя телефон. Материал имеет название количество цену. Объект содержит этапы. Подрядчик выполняет этапы. Этап использует материалы.",
			entities:  []string{"Объект", "Этап", "Подрядчик", "Материал"},
			fragments: []string{"CREATE DATABASE stroitelstvo", "byudzhet NUMERIC(12,2)", "data_nachalo DATE", "data_okonchanie DATE", "obekt_id INTEGER", "podryadchik_id INTEGER", "material_id INTEGER"},
		},
		{
			name:      "service-center",
			input:     "База данных сервиса. Клиент имеет имя телефон email. Устройство имеет тип модель серийный номер гарантию. Заявка имеет дату описание статус стоимость. Мастер имеет имя телефон. Клиент создает заявки. Заявка связана с устройством. Мастер ремонтирует устройства.",
			entities:  []string{"Клиент", "Устройство", "Заявка", "Мастер"},
			fragments: []string{"CREATE DATABASE servis", "tip VARCHAR(255)", "model VARCHAR(255)", "garantiya INTEGER", "klient_id INTEGER", "ustroystvo_id INTEGER"},
		},
		{
			name:      "tourism",
			input:     "Система туризма. Турист имеет имя телефон паспорт. Тур имеет название страну маршрут цену дату начала дату окончания. Бронирование имеет дату статус сумму. Менеджер имеет имя email. Турист бронирует туры. Бронирование связано с туром. Менеджер ведет бронирования.",
			entities:  []string{"Турист", "Тур", "Бронирование", "Менеджер"},
			fragments: []string{"CREATE DATABASE turizm", "strana VARCHAR(255)", "marshrut VARCHAR(255)", "data_nachalo DATE", "data_okonchanie DATE", "turist_id INTEGER", "tur_id INTEGER"},
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
