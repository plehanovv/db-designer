package nlp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

func ProcessText(text string) Document {
	requestBody := AnalyzeRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Document{}
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Post(
		nlpServiceURL(),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return ProcessTextLocally(text)
	}

	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return ProcessTextLocally(text)
	}

	var document Document

	err = json.NewDecoder(response.Body).Decode(&document)
	if err != nil {
		return ProcessTextLocally(text)
	}

	if len(document.Tokens) == 0 {
		return ProcessTextLocally(text)
	}

	if !hasUsefulPOSTags(document) {
		return ProcessTextLocally(text)
	}

	document.Source = "spacy"
	return document
}

func hasUsefulPOSTags(document Document) bool {
	for _, token := range document.Tokens {
		switch token.Pos {
		case "NOUN", "PROPN", "VERB":
			return true
		}
	}

	return false
}

func nlpServiceURL() string {
	value := os.Getenv("NLP_SERVICE_URL")
	if value == "" {
		return "http://localhost:8000/analyze"
	}

	return value
}

func ProcessTextLocally(text string) Document {
	tokens := make([]Token, 0)
	var builder strings.Builder
	sentence := 0

	flush := func() {
		if builder.Len() == 0 {
			return
		}

		word := builder.String()
		lemma := strings.ToLower(word)
		tokens = append(tokens, Token{
			Text:      word,
			Lemma:     lemma,
			Pos:       detectLocalPOS(lemma),
			Index:     len(tokens),
			HeadIndex: len(tokens),
			Sentence:  sentence,
		})
		builder.Reset()
	}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			continue
		}

		flush()
		if r == '.' || r == '!' || r == '?' || r == ';' || r == '\n' {
			sentence++
		}
	}
	flush()

	return Document{
		Tokens:   tokens,
		Language: detectLocalLanguage(text),
		Source:   "local_rule_fallback",
	}
}

func detectLocalLanguage(text string) string {
	for _, r := range text {
		if r >= 'а' && r <= 'я' || r == 'ё' || r == 'Ё' {
			return "ru"
		}
	}

	return "en"
}

func detectLocalPOS(lemma string) string {
	if localRelationWords[lemma] {
		return "VERB"
	}
	if localNoiseWords[lemma] {
		return "X"
	}
	return "NOUN"
}

var localRelationWords = termSet(
	// Common relation verbs used in controlled English examples.
	"have", "has", "contain", "contains", "include",
	"includes", "belong", "belongs", "connect", "connects",
	"store", "stores", "create", "creates", "buy",
	"buys", "borrow", "borrows", "pay", "pays",
	"use", "uses", "assign", "assigns", "enroll",
	"enrolls", "teach", "teaches", "manage", "manages",
	"register", "registers", "book", "books", "reserve",
	"reserves", "deliver", "delivers", "receive", "receives",
	"supply", "supplies", "sell", "sells", "publish",
	"publishes", "write", "writes", "submit", "submits",
	"pass", "passes", "evaluate", "evaluates",

	// Russian verbs for the main demo domains: library, study process, commerce and logistics.
	"имеет", "имеют", "иметь", "есть", "может",
	"содержит", "содержат", "включает", "принадлежит", "связан",
	"связана", "связано", "связать", "относится", "хранит",
	"хранится", "хранятся", "оформляет", "оформляют", "создает",
	"создают", "создавать", "обрабатывает", "обрабатывают", "заказывает",
	"заказывают", "покупает", "покупают", "оплачивает", "оплачивают",
	"получает", "получают", "получать", "поставляет", "поставляют",
	"поставлять", "доставляет", "доставляют", "доставлять", "выдает",
	"выдают", "выдавать", "берет", "берут", "брать",
	"бронирует", "бронируют", "бронировать", "резервирует", "резервировать",
	"работает", "работают", "ведет", "ведут", "управляет",
	"управляют", "выполняет", "выполняют", "выполнять", "участвует",
	"участвуют", "регистрирует", "записывается", "проходит", "проходят", "сдает",
	"сдают", "сдавать", "оценивает", "оценивают", "оценивать",
	"посещает", "посещают", "посещать", "публикует", "публикуют",
	"публиковать", "пишет", "пишут", "писать",
)

var localNoiseWords = termSet(
	"a", "an", "and", "are", "by",
	"each", "for", "from", "in", "is",
	"of", "the", "to", "with",
	"в", "а", "и", "или", "также",
	"еще", "ещё", "тоже", "плюс", "для",
	"на", "с", "со", "по", "у",
	"который", "которая", "которые", "которых", "каждый",
	"каждая", "каждое",
)

func termSet(values ...string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}
