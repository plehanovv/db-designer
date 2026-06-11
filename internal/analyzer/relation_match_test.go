package analyzer

import (
	"testing"

	"db-designer-vkr/internal/nlp"
)

func TestRelationMatchesUseDependencySubjectPredicateObject(t *testing.T) {
	document := nlp.Document{
		Tokens: []nlp.Token{
			{Text: "Order", Lemma: "order", Pos: "NOUN", Dependency: "nsubj", Head: "belongs", Index: 0, HeadIndex: 1, Sentence: 0},
			{Text: "belongs", Lemma: "belong", Pos: "VERB", Dependency: "ROOT", Head: "belongs", Index: 1, HeadIndex: 1, Sentence: 0},
			{Text: "to", Lemma: "to", Pos: "ADP", Dependency: "prep", Head: "belongs", Index: 2, HeadIndex: 1, Sentence: 0},
			{Text: "Customer", Lemma: "customer", Pos: "NOUN", Dependency: "pobj", Head: "to", Index: 3, HeadIndex: 2, Sentence: 0},
		},
	}

	matches := RelationMatches(document)
	if len(matches) != 1 {
		t.Fatalf("expected 1 dependency relation match, got %d: %#v", len(matches), matches)
	}

	match := matches[0]
	if normalizeEntity(match.From.Lemma) != "Order" || normalizeEntity(match.To.Lemma) != "Customer" {
		t.Fatalf("unexpected relation endpoints: %#v", match)
	}

	if match.Rule != "dependency_subject_predicate_object" {
		t.Fatalf("expected dependency rule, got %q", match.Rule)
	}
}

func TestRelationMatchesUseRelativeAssociatedClause(t *testing.T) {
	document := nlp.Document{
		Tokens: []nlp.Token{
			{Text: "История", Lemma: "история", Pos: "NOUN", Index: 0, Sentence: 0},
			{Text: "заказов", Lemma: "заказ", Pos: "NOUN", Index: 1, Sentence: 0},
			{Text: "которая", Lemma: "который", Pos: "PRON", Index: 2, Sentence: 0},
			{Text: "связана", Lemma: "связать", Pos: "VERB", Index: 3, Sentence: 0},
			{Text: "с", Lemma: "с", Pos: "ADP", Index: 4, Sentence: 0},
			{Text: "пользователями", Lemma: "пользователь", Pos: "NOUN", Index: 5, Sentence: 0},
		},
	}

	matches := RelationMatches(document)
	if len(matches) != 1 {
		t.Fatalf("expected 1 relative relation match, got %d: %#v", len(matches), matches)
	}

	match := matches[0]
	if normalizeEntity(match.From.Lemma) != "История" || normalizeEntity(match.To.Lemma) != "Пользователь" {
		t.Fatalf("unexpected relation endpoints: %#v", match)
	}

	if match.Type != "associated_with" {
		t.Fatalf("expected associated relation, got %q", match.Type)
	}
}

func TestRelationMatchesUseRussianNLPLemmas(t *testing.T) {
	document := nlp.Document{
		Tokens: []nlp.Token{
			{Text: "Студент", Lemma: "студент", Pos: "NOUN", Index: 0, Sentence: 0},
			{Text: "записывается", Lemma: "записываться", Pos: "VERB", Index: 1, Sentence: 0},
			{Text: "на", Lemma: "на", Pos: "ADP", Index: 2, Sentence: 0},
			{Text: "курсы", Lemma: "курс", Pos: "NOUN", Index: 3, Sentence: 0},

			{Text: "Преподаватель", Lemma: "преподаватель", Pos: "NOUN", Index: 4, Sentence: 1},
			{Text: "ведет", Lemma: "вести", Pos: "VERB", Index: 5, Sentence: 1},
			{Text: "курсы", Lemma: "курс", Pos: "NOUN", Index: 6, Sentence: 1},

			{Text: "Студент", Lemma: "студент", Pos: "NOUN", Index: 7, Sentence: 2},
			{Text: "сдает", Lemma: "сдавать", Pos: "VERB", Index: 8, Sentence: 2},
			{Text: "экзамены", Lemma: "экзамен", Pos: "NOUN", Index: 9, Sentence: 2},

			{Text: "Преподаватель", Lemma: "преподаватель", Pos: "NOUN", Index: 10, Sentence: 3},
			{Text: "оценивает", Lemma: "оценивать", Pos: "VERB", Index: 11, Sentence: 3},
			{Text: "результаты", Lemma: "результат", Pos: "NOUN", Index: 12, Sentence: 3},
		},
	}

	matches := RelationMatches(document)
	for _, expected := range []struct {
		from string
		to   string
	}{
		{from: "Студент", to: "Курс"},
		{from: "Преподаватель", to: "Курс"},
		{from: "Студент", to: "Экзамен"},
		{from: "Преподаватель", to: "Результат"},
	} {
		if !hasRelationMatch(matches, expected.from, expected.to) {
			t.Fatalf("expected relation %s -> %s, got %#v", expected.from, expected.to, matches)
		}
	}
}

func hasRelationMatch(matches []RelationMatch, from string, to string) bool {
	for _, match := range matches {
		if normalizeEntity(match.From.Lemma) == from && normalizeEntity(match.To.Lemma) == to {
			return true
		}
	}

	return false
}
