package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"db-designer-vkr/internal/model"
)

func ExtractDatabase(text string) model.Database {
	normalized := normalizeSpaces(strings.ToLower(text))
	if normalized == "" {
		return model.Database{Name: "database"}
	}

	if name, ok := databaseNameFromRussianIntro(normalized); ok {
		return model.Database{Name: titleWords(name), Domain: titleWords(name)}
	}

	if name, ok := databaseNameFromEnglishIntro(normalized); ok {
		return model.Database{Name: titleWords(name), Domain: titleWords(name)}
	}

	return model.Database{Name: "database"}
}

func databaseNameFromRussianIntro(text string) (string, bool) {
	markers := []string{
		"база данных",
		"база для",
		"бд",
		"предметная область",
		"предметной области",
		"система",
	}

	for _, marker := range markers {
		index := strings.Index(text, marker)
		if index < 0 {
			continue
		}

		rest := strings.TrimSpace(text[index+len(marker):])
		rest = strings.TrimLeft(rest, " :-—")
		candidate := cleanDatabaseName(firstClause(rest))
		if candidate != "" {
			return candidate, true
		}
	}

	return "", false
}

func databaseNameFromEnglishIntro(text string) (string, bool) {
	if strings.Contains(text, "online shop") || strings.Contains(text, "internet shop") {
		return "online shop", true
	}

	markers := []string{
		"database for",
		"database of",
		"domain",
		"system",
	}

	for _, marker := range markers {
		index := strings.Index(text, marker)
		if index < 0 {
			continue
		}

		rest := strings.TrimSpace(text[index+len(marker):])
		rest = strings.TrimLeft(rest, " :-—")
		candidate := cleanDatabaseName(firstClause(rest))
		if candidate != "" {
			return candidate, true
		}
	}

	return "", false
}

func firstClause(value string) string {
	for index, r := range value {
		if r == '.' || r == '!' || r == '?' || r == ';' || r == '\n' {
			return value[:index]
		}
	}

	return value
}

func cleanDatabaseName(value string) string {
	words := strings.Fields(value)
	cleaned := make([]string, 0, len(words))

	for _, word := range words {
		word = normalizeWord(word)
		if word == "" || isDatabaseNameNoise(word) {
			continue
		}
		if override, exists := databaseWordOverrides[word]; exists {
			word = override
		} else if !isASCII(word) {
			word = singularRussian(word)
		}

		cleaned = append(cleaned, word)
		if len(cleaned) == 3 {
			break
		}
	}

	return strings.Join(cleaned, " ")
}

func isDatabaseNameNoise(value string) bool {
	noise := map[string]bool{
		"для":     true,
		"по":      true,
		"система": true,
		"системы": true,
		"the":     true,
		"a":       true,
		"an":      true,
	}

	return noise[normalizeWord(value)]
}

var databaseWordOverrides = map[string]string{
	"библиотеки":       "библиотека",
	"магазина":         "магазин",
	"магазинов":        "магазин",
	"учета":            "учет",
	"заявок":           "заявка",
	"гостиницы":        "гостиница",
	"клиники":          "клиника",
	"аренды":           "аренда",
	"автомобилей":      "автомобиль",
	"управления":       "управление",
	"проектами":        "проект",
	"проекта":          "проект",
	"агентства":        "агентство",
	"недвижимости":     "недвижимость",
	"логистики":        "логистика",
	"банка":            "банк",
	"мероприятий":      "мероприятие",
	"подбора":          "подбор",
	"персонала":        "персонал",
	"ресторана":        "ресторан",
	"производства":     "производство",
	"сайта":            "сайт",
	"страхования":      "страхование",
	"авиаперевозок":    "авиаперевозка",
	"фитнеса":          "фитнес",
	"экзаменов":        "экзамен",
	"документооборота": "документооборот",
	"аптеки":           "аптека",
	"кинотеатра":       "кинотеатр",
	"музея":            "музей",
	"строительства":    "строительство",
	"сервиса":          "сервис",
	"туризма":          "туризм",
}

func normalizeSpaces(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func titleWords(value string) string {
	words := strings.Fields(value)
	for i, word := range words {
		first, size := utf8.DecodeRuneInString(word)
		if first == utf8.RuneError {
			continue
		}
		words[i] = string(unicode.ToUpper(first)) + word[size:]
	}

	return strings.Join(words, " ")
}
