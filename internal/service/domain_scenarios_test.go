package service

import (
	"strings"
	"testing"
)

func TestAdditionalDomainScenarios(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	scenarios := []struct {
		name      string
		input     string
		entities  []string
		fragments []string
	}{
		{
			name:  "warehouse",
			input: "База данных склада. Поставщик имеет имя телефон. Товар имеет название цену количество. Поставка имеет дату сумму. Поставщик создает поставки. Поставка содержит товары.",
			entities: []string{
				"Поставщик", "Товар", "Поставка",
			},
			fragments: []string{"CREATE DATABASE sklad", "postavschik_id INTEGER", "postavka_id INTEGER", "CREATE INDEX idx_postavka_postavschik_id"},
		},
		{
			name:  "crm",
			input: "База данных CRM. Клиент имеет имя телефон почту. Менеджер имеет имя email. Сделка имеет сумму статус дату создания. Менеджер обрабатывает сделки. Сделка связана с клиентом.",
			entities: []string{
				"Клиент", "Менеджер", "Сделка",
			},
			fragments: []string{"CREATE DATABASE crm", "summa NUMERIC(12,2)", "data_sozdanie DATE", "menedzher_id INTEGER", "klient_id INTEGER"},
		},
		{
			name:  "project-management",
			input: "Система управления проектами. Проект имеет название код. Задача имеет название статус дату окончания. Пользователь создает задачи. Задача относится к проекту.",
			entities: []string{
				"Проект", "Задача", "Пользователь",
			},
			fragments: []string{"CREATE DATABASE upravlenie_proekt", "data_okonchanie DATE", "polzovatel_id INTEGER", "proekt_id INTEGER"},
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
