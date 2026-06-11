package analyzer

import (
	"testing"

	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func TestExtractAttributesFromSpacyRelativePossessiveClause(t *testing.T) {
	document := nlp.Document{
		Tokens: []nlp.Token{
			{Text: "Пользователи", Lemma: "пользователь", Pos: "NOUN", Index: 0, Sentence: 0},
			{Text: "у", Lemma: "у", Pos: "ADP", Index: 1, Sentence: 0},
			{Text: "которых", Lemma: "который", Pos: "PRON", Index: 2, Sentence: 0},
			{Text: "есть", Lemma: "быть", Pos: "AUX", Index: 3, Sentence: 0},
			{Text: "номер", Lemma: "номер", Pos: "NOUN", Index: 4, Sentence: 0},
			{Text: "телефона", Lemma: "телефон", Pos: "NOUN", Index: 5, Sentence: 0},
			{Text: "почта", Lemma: "почта", Pos: "NOUN", Index: 6, Sentence: 0},
			{Text: "дата", Lemma: "дата", Pos: "NOUN", Index: 7, Sentence: 0},
			{Text: "последней", Lemma: "последний", Pos: "ADJ", Index: 8, Sentence: 0},
			{Text: "покупки", Lemma: "покупка", Pos: "NOUN", Index: 9, Sentence: 0},
		},
	}

	entityMap := map[string]*model.Entity{
		"Пользователь": {Name: "Пользователь"},
	}

	ExtractAttributes(document, entityMap)

	entity := entityMap["Пользователь"]
	for _, name := range []string{"телефон", "почта", "дата_покупка"} {
		if !testHasAttribute(entity.Attributes, name) {
			t.Fatalf("expected attribute %q, got %#v", name, entity.Attributes)
		}
	}

	for _, name := range []string{"номер", "последний", "покупка"} {
		if testHasAttribute(entity.Attributes, name) {
			t.Fatalf("attribute %q must be skipped, got %#v", name, entity.Attributes)
		}
	}
}

func TestExtractAttributesSkipsStockLocationPhrase(t *testing.T) {
	document := nlp.Document{
		Tokens: []nlp.Token{
			{Text: "товары", Lemma: "товар", Pos: "NOUN", Index: 0, Sentence: 0},
			{Text: "которые", Lemma: "который", Pos: "PRON", Index: 1, Sentence: 0},
			{Text: "имеют", Lemma: "иметь", Pos: "VERB", Index: 2, Sentence: 0},
			{Text: "цену", Lemma: "цена", Pos: "NOUN", Index: 3, Sentence: 0},
			{Text: "скидку", Lemma: "скидка", Pos: "NOUN", Index: 4, Sentence: 0},
			{Text: "количество", Lemma: "количество", Pos: "NOUN", Index: 5, Sentence: 0},
			{Text: "на", Lemma: "на", Pos: "ADP", Index: 6, Sentence: 0},
			{Text: "складе", Lemma: "склад", Pos: "NOUN", Index: 7, Sentence: 0},
		},
	}

	entityMap := map[string]*model.Entity{
		"Товар": {Name: "Товар"},
	}

	ExtractAttributes(document, entityMap)

	entity := entityMap["Товар"]
	for _, name := range []string{"цена", "скидка", "количество"} {
		if !testHasAttribute(entity.Attributes, name) {
			t.Fatalf("expected attribute %q, got %#v", name, entity.Attributes)
		}
	}

	if testHasAttribute(entity.Attributes, "склад") {
		t.Fatalf("stock location must not be emitted as scalar attribute, got %#v", entity.Attributes)
	}
}

func testHasAttribute(attributes []model.Attribute, name string) bool {
	for _, attribute := range attributes {
		if attribute.Name == name {
			return true
		}
	}

	return false
}
