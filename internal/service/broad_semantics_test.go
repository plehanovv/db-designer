package service

import (
	"strings"
	"testing"
)

func TestBroadDomainSemantics(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := []struct {
		name      string
		input     string
		entities  []string
		fragments []string
	}{
		{
			name:      "real-estate",
			input:     "База данных агентства недвижимости. Клиент имеет имя телефон email. Объект имеет адрес город цену площадь статус. Агент имеет имя телефон. Сделка имеет сумму дату статус. Агент ведет сделки. Сделка связана с клиентом. Сделка связана с объектом.",
			entities:  []string{"Клиент", "Объект", "Агент", "Сделка"},
			fragments: []string{"CREATE DATABASE agentstvo_nedvizhimost", "adres VARCHAR(255)", "cena NUMERIC(12,2)", "agent_id INTEGER", "klient_id INTEGER", "obekt_id INTEGER"},
		},
		{
			name:      "logistics",
			input:     "Система логистики. Клиент имеет имя телефон адрес. Груз имеет название вес объем статус. Курьер имеет имя телефон. Доставка имеет дату стоимость статус. Клиент создает доставки. Курьер доставляет грузы. Доставка связана с грузом.",
			entities:  []string{"Клиент", "Груз", "Курьер", "Доставка"},
			fragments: []string{"CREATE DATABASE logistika", "ves NUMERIC(12,2)", "obem NUMERIC(12,2)", "stoimost NUMERIC(12,2)", "klient_id INTEGER", "gruz_id INTEGER"},
		},
		{
			name:      "finance",
			input:     "База данных банка. Клиент имеет имя телефон паспорт. Счет имеет номер баланс статус. Карта имеет номер срок статус. Операция имеет сумму дату тип. Клиент открывает счета. Счет содержит операции. Карта связана со счетом.",
			entities:  []string{"Клиент", "Счет", "Карта", "Операция"},
			fragments: []string{"CREATE DATABASE bank", "balans NUMERIC(12,2)", "summa NUMERIC(12,2)", "schet_id INTEGER", "klient_id INTEGER"},
		},
		{
			name:      "events",
			input:     "Система мероприятий. Участник имеет имя email телефон. Событие имеет название дату адрес статус. Билет имеет номер цену статус. Участник покупает билеты. Билет связан с событием.",
			entities:  []string{"Участник", "Событие", "Билет"},
			fragments: []string{"CREATE DATABASE meropriyatie", "data DATE", "adres VARCHAR(255)", "cena NUMERIC(12,2)", "uchastnik_id INTEGER", "sobytie_id INTEGER"},
		},
		{
			name:      "hr",
			input:     "База данных подбора персонала. Кандидат имеет имя телефон email резюме. Вакансия имеет название зарплату статус. Собеседование имеет дату время оценку. Кандидат откликается на вакансии. Собеседование связано с кандидатом. Собеседование связано с вакансией.",
			entities:  []string{"Кандидат", "Вакансия", "Собеседование"},
			fragments: []string{"CREATE DATABASE podbor_personal", "zarplata NUMERIC(12,2)", "vremya TIME", "ocenka INTEGER", "kandidat_id INTEGER", "vakansiya_id INTEGER"},
		},
		{
			name:      "restaurant",
			input:     "База данных ресторана. Гость имеет имя телефон. Стол имеет номер количество мест статус. Бронирование имеет дату время статус. Заказ имеет сумму статус. Гость бронирует столы. Бронирование связано со столом. Гость оформляет заказы. Заказ содержит блюда.",
			entities:  []string{"Гость", "Стол", "Бронирование", "Заказ", "Блюдо"},
			fragments: []string{"CREATE DATABASE restoran", "kolichestvo INTEGER", "vremya TIME", "gost_id INTEGER", "stol_id INTEGER", "zakaz_id INTEGER"},
		},
		{
			name:      "manufacturing",
			input:     "Система производства. Изделие имеет название код стоимость. Материал имеет название количество цену. Операция имеет название длительность статус. Заказ содержит изделия. Изделие использует материалы. Операция связана с изделием.",
			entities:  []string{"Изделие", "Материал", "Операция", "Заказ"},
			fragments: []string{"CREATE DATABASE proizvodstvo", "stoimost NUMERIC(12,2)", "dlitelnost INTEGER", "zakaz_id INTEGER", "izdelie_id INTEGER"},
		},
		{
			name:      "content",
			input:     "База данных сайта. Автор имеет имя email роль. Статья имеет название текст дату статус. Комментарий имеет текст дату статус. Автор публикует статьи. Статья содержит комментарии.",
			entities:  []string{"Автор", "Статья", "Комментарий"},
			fragments: []string{"CREATE DATABASE sayt", "rol VARCHAR(255)", "data DATE", "avtor_id INTEGER", "statya_id INTEGER"},
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
