package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var ignoredEntityLemmas = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "by": true, "each": true,
	"for": true, "from": true, "in": true, "is": true, "of": true, "the": true,
	"to": true, "with": true, "one": true, "many": true, "multiple": true,
	"several": true, "may": true, "can": true,
	"amount": true, "date": true, "email": true, "name": true, "number": true,
	"phone": true, "title": true, "required": true, "mandatory": true, "unique": true,
	"isbn": true, "place": true, "places": true,
	"мне":              true,
	"ней":              true,
	"нужен":            true,
	"нужна":            true,
	"нужно":            true,
	"а":                true,
	"также":            true,
	"еще":              true,
	"ещё":              true,
	"тоже":             true,
	"плюс":             true,
	"атрибут":          true,
	"атрибуты":         true,
	"база":             true,
	"бд":               true,
	"библиотек":        true,
	"библиотека":       true,
	"библиотеки":       true,
	"данн":             true,
	"данных":           true,
	"данные":           true,
	"дата":             true,
	"время":            true,
	"интернет":         true,
	"запись":           true,
	"имя":              true,
	"инн":              true,
	"информация":       true,
	"каждый":           true,
	"каждая":           true,
	"каждое":           true,
	"код":              true,
	"количество":       true,
	"которая":          true,
	"который":          true,
	"которые":          true,
	"которых":          true,
	"магазин":          true,
	"магазина":         true,
	"может":            true,
	"нас":              true,
	"один":             true,
	"одна":             true,
	"одно":             true,
	"несколько":        true,
	"область":          true,
	"области":          true,
	"окончание":        true,
	"окончания":        true,
	"обязательный":     true,
	"обязательная":     true,
	"обязательное":     true,
	"описание":         true,
	"поле":             true,
	"поля":             true,
	"почта":            true,
	"покупка":          true,
	"покупки":          true,
	"начала":           true,
	"последний":        true,
	"последней":        true,
	"рождение":         true,
	"рождения":         true,
	"предметная":       true,
	"режим":            true,
	"работ":            true,
	"работа":           true,
	"система":          true,
	"систем":           true,
	"системе":          true,
	"аренд":            true,
	"аренда":           true,
	"аренды":           true,
	"гостиниц":         true,
	"гостиница":        true,
	"гостиницы":        true,
	"клиник":           true,
	"клиника":          true,
	"клиники":          true,
	"скидка":           true,
	"структура":        true,
	"сумма":            true,
	"таблица":          true,
	"таблицы":          true,
	"телефон":          true,
	"фио":              true,
	"вместимость":      true,
	"вес":              true,
	"госномер":         true,
	"грузоподъемность": true,
	"грузоподъёмность": true,
	"марка":            true,
	"модель":           true,
	"объем":            true,
	"объём":            true,
	"площадь":          true,
	"стоимость":        true,
	"управление":       true,
	"управления":       true,
	"учебн":            true,
	"процесса":         true,
	"уникальный":       true,
	"уникальная":       true,
	"уникальное":       true,
	"хранение":         true,
	"хранения":         true,
	"хранится":         true,
	"связан":           true,
	"связана":          true,
	"связано":          true,
	"экземпляр":        true,
	"цена":             true,
	"агентство":        true,
	"агентства":        true,
	"банк":             true,
	"банка":            true,
	"логистик":         true,
	"логистика":        true,
	"логистики":        true,
	"складская":        true,
	"недвижимость":     true,
	"недвижимости":     true,
	"персонал":         true,
	"персонала":        true,
	"подбор":           true,
	"подбора":          true,
	"производство":     true,
	"производства":     true,
	"ресторан":         true,
	"ресторана":        true,
	"сайт":             true,
	"сайта":            true,
	"страхование":      true,
	"страхования":      true,
	"страхов":          true,
	"страховой":        true,
	"страховые":        true,
	"страховым":        true,
	"авиаперевозка":    true,
	"авиаперевозок":    true,
	"фитнес":           true,
	"фитнеса":          true,
	"документооборот":  true,
	"документооборота": true,
	"аптека":           true,
	"аптеки":           true,
	"кинотеатр":        true,
	"кинотеатра":       true,
	"музей":            true,
	"музея":            true,
	"строительство":    true,
	"строительства":    true,
	"сервис":           true,
	"сервиса":          true,
	"туризм":           true,
	"туризма":          true,
	"центр":            true,
	"центра":           true,
	"центром":          true,
	"возникает":        true,
	"возникают":        true,
	"выполняется":      true,
	"логистическ":      true,
	"описывает":        true,
	"размещается":      true,
	"распределительн":  true,
	"транспортн":       true,
	"транспортно":      true,
	"транспортные":     true,
}

var ignoredAttributeLemmas = map[string]bool{
	"a": true, "an": true, "and": true, "or": true, "the": true,
	"with": true, "of": true, "to": true, "required": true, "mandatory": true, "unique": true,
	"а": true, "и": true, "или": true, "на": true,
	"также": true, "еще": true,
	"ещё": true, "тоже": true,
	"плюс":         true,
	"обязательный": true,
	"обязательная": true,
	"обязательное": true,
	"последний":    true,
	"последней":    true,
	"выезд":        true,
	"выезда":       true,
	"выдача":       true,
	"выдачи":       true,
	"заезд":        true,
	"заезда":       true,
	"рождение":     true,
	"рождения":     true,
	"создание":     true,
	"создания":     true,
	"возврат":      true,
	"возврата":     true,
	"начало":       true,
	"начала":       true,
	"окончание":    true,
	"окончания":    true,
	"покупка":      true,
	"покупки":      true,
	"с":            true, "склад": true,
	"складе": true,
	"со":     true, "в": true, "у": true,
	"уникальный": true,
	"уникальная": true,
	"уникальное": true,
}

var russianSingularReplacements = []struct {
	suffix string
	value  string
}{
	{"иями", "ия"},
	{"ями", "я"},
	{"ами", "а"},
	{"ого", ""},
	{"ему", ""},
	{"ыми", ""},
	{"ими", ""},
	{"ые", ""},
	{"ие", "ие"},
	{"ов", ""},
	{"ев", ""},
	{"ей", "ь"},
	{"ам", "а"},
	{"ям", "я"},
	{"ах", "а"},
	{"ях", "я"},
	{"ом", ""},
	{"ем", "ь"},
	{"ой", "а"},
	{"ых", ""},
	{"их", ""},
	{"у", ""},
	{"е", ""},
	{"ы", ""},
	{"и", ""},
}

var normalizedWordOverrides = map[string]string{
	"автора":       "автор",
	"группу":       "группа",
	"дату":         "дата",
	"должность":    "должность",
	"выезда":       "выезд",
	"выдачи":       "выдача",
	"заезда":       "заезд",
	"рождения":     "рождение",
	"создания":     "создание",
	"возврата":     "возврат",
	"начала":       "начало",
	"окончания":    "окончание",
	"тему":         "тема",
	"роль":         "роль",
	"марку":        "марка",
	"которая":      "который",
	"которые":      "который",
	"которых":      "который",
	"первой":       "первый",
	"первая":       "первый",
	"покупки":      "покупка",
	"последней":    "последний",
	"последняя":    "последний",
	"обязательную": "обязательный",
	"уникальную":   "уникальный",
	"сумму":        "сумма",
	"скидку":       "скидка",
	"телефона":     "телефон",
	"цену":         "цена",
	"адреса":       "адрес",
	"баланса":      "баланс",
	"бюджета":      "бюджет",
	"веса":         "вес",
	"города":       "город",
	"категорию":    "категория",
	"комментарий":  "комментарий",
	"оплату":       "платеж",
	"остатка":      "остаток",
	"оценку":       "оценка",
	"площадь":      "площадь",
	"приоритета":   "приоритет",
	"срока":        "срок",
	"уровня":       "уровень",
	"зарплату":     "зарплата",
	"премию":       "премия",
	"выплату":      "выплата",
	"класса":       "класс",
	"кабинета":     "кабинет",
	"места":        "место",
	"дозировку":    "дозировка",
	"жанра":        "жанр",
	"серию":        "серия",
	"гарантию":     "гарантия",
	"маршрута":     "маршрут",
	"страну":       "страна",
}

var normalizedEntityOverrides = map[string]string{
	"categories":     "category",
	"companies":      "company",
	"deliveries":     "delivery",
	"авторы":         "автор",
	"читатели":       "читатель",
	"читателей":      "читатель",
	"заказов":        "заказ",
	"заказа":         "заказ",
	"заказы":         "заказ",
	"заявки":         "заявка",
	"заявок":         "заявка",
	"категории":      "категория",
	"категорию":      "категория",
	"кафедры":        "кафедра",
	"книгой":         "книга",
	"книги":          "книга",
	"курсы":          "курс",
	"платежа":        "платеж",
	"платежи":        "платеж",
	"поставки":       "поставка",
	"проекта":        "проект",
	"проекты":        "проект",
	"сделки":         "сделка",
	"задачи":         "задача",
	"прием":          "прием",
	"клиента":        "клиент",
	"клиенту":        "клиент",
	"клиентов":       "клиент",
	"бронирования":   "бронирование",
	"приемы":         "прием",
	"приемов":        "прием",
	"врачу":          "врач",
	"сотрудники":     "сотрудник",
	"сотруднику":     "сотрудник",
	"автомобили":     "автомобиль",
	"автомобилей":    "автомобиль",
	"пользователи":   "пользователь",
	"пользователей":  "пользователь",
	"пользователями": "пользователь",
	"товары":         "товар",
	"товаров":        "товар",
	"абонементы":     "абонемент",
	"адреса":         "адрес",
	"билеты":         "билет",
	"вакансии":       "вакансия",
	"выплаты":        "выплата",
	"выдач":          "выдача",
	"грузы":          "груз",
	"груза":          "груз",
	"группе":         "группа",
	"группы":         "группа",
	"документы":      "документ",
	"доставки":       "доставка",
	"здания":         "здание",
	"занятия":        "занятие",
	"комнаты":        "комната",
	"контракты":      "контракт",
	"курьеры":        "курьер",
	"маршрута":       "маршрут",
	"маршруты":       "маршрут",
	"материалы":      "материал",
	"операции":       "операция",
	"отзывы":         "отзыв",
	"оплаты":         "оплата",
	"оплат":          "оплата",
	"рейсы":          "рейс",
	"рейса":          "рейс",
	"резюме":         "резюме",
	"склады":         "склад",
	"складе":         "склад",
	"склада":         "склад",
	"события":        "событие",
	"статьи":         "статья",
	"участники":      "участник",
	"филиалы":        "филиал",
	"экзамены":       "экзамен",
	"блюда":          "блюдо",
	"вакансией":      "вакансия",
	"вакансию":       "вакансия",
	"изделием":       "изделие",
	"изделия":        "изделие",
	"комментарии":    "комментарий",
	"событием":       "событие",
	"счета":          "счет",
	"счетом":         "счет",
	"полисы":         "полис",
	"полисов":        "полис",
	"полисом":        "полис",
	"случаи":         "случай",
	"случаев":        "случай",
	"случаем":        "случай",
	"пассажиры":      "пассажир",
	"пассажиров":     "пассажир",
	"аэропорты":      "аэропорт",
	"аэропортом":     "аэропорт",
	"тренеры":        "тренер",
	"тренером":       "тренер",
	"залы":           "зал",
	"залом":          "зал",
	"ученики":        "ученик",
	"учеников":       "ученик",
	"результаты":     "результат",
	"результатом":    "результат",
	"согласования":   "согласование",
	"подписи":        "подпись",
	"исполнители":    "исполнитель",
	"исполнителю":    "исполнитель",
	"лекарства":      "лекарство",
	"препараты":      "препарат",
	"препаратов":     "препарат",
	"продажи":        "продажа",
	"сеансы":         "сеанс",
	"сеансом":        "сеанс",
	"фильмы":         "фильм",
	"фильмом":        "фильм",
	"экспонаты":      "экспонат",
	"экспонатов":     "экспонат",
	"выставки":       "выставка",
	"этапы":          "этап",
	"этапом":         "этап",
	"устройства":     "устройство",
	"устройством":    "устройство",
	"мастера":        "мастер",
	"мастеру":        "мастер",
	"туры":           "тур",
	"туров":          "тур",
	"туристы":        "турист",
	"туристов":       "турист",
	"водители":       "водитель",
	"водителя":       "водитель",
	"зоны":           "зона",
	"зоне":           "зона",
	"инциденты":      "инцидент",
	"накладные":      "накладная",
	"накладна":       "накладная",
	"накладной":      "накладная",
	"пункты":         "пункт",
	"пункта":         "пункт",
	"средства":       "транспортное средство",
	"средство":       "транспортное средство",
	"ячейки":         "ячейка",
	"ячейке":         "ячейка",
	"ячеек":          "ячейка",
}

func normalizeEntity(value string) string {
	value = normalizeWord(value)
	if value == "" {
		return value
	}

	overridden := false
	if override, exists := normalizedEntityOverrides[value]; exists {
		value = override
		overridden = true
	}

	if isASCII(value) && len(value) > 3 && strings.HasSuffix(value, "s") {
		value = strings.TrimSuffix(value, "s")
	}

	if !isASCII(value) && !overridden {
		value = singularRussian(value)
		if override, exists := normalizedEntityOverrides[value]; exists {
			value = override
		}
	}

	first, size := utf8.DecodeRuneInString(value)
	return string(unicode.ToUpper(first)) + value[size:]
}

func normalizeAttribute(value string) string {
	return normalizeDomainWord(value)
}

func normalizeWord(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Trim(value, ".,;:!?()[]{}\"'`")
	return value
}

func isIgnoredEntity(value string) bool {
	value = normalizeDomainWord(value)
	if ignoredEntityLemmas[value] {
		return true
	}

	if override, exists := normalizedEntityOverrides[value]; exists && ignoredEntityLemmas[override] {
		return true
	}

	if !isASCII(value) && ignoredEntityLemmas[singularRussian(value)] {
		return true
	}

	return false
}

func isIgnoredAttribute(value string) bool {
	return ignoredAttributeLemmas[normalizeDomainWord(value)]
}

func normalizeDomainWord(value string) string {
	value = normalizeWord(value)
	if override, exists := normalizedWordOverrides[value]; exists {
		return override
	}

	return value
}

func singularRussian(value string) string {
	for _, replacement := range russianSingularReplacements {
		if len([]rune(value)) > len([]rune(replacement.suffix))+2 && strings.HasSuffix(value, replacement.suffix) {
			return strings.TrimSuffix(value, replacement.suffix) + replacement.value
		}
	}
	return value
}

func isASCII(value string) bool {
	for _, r := range value {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func startsWithUpper(value string) bool {
	first, _ := utf8.DecodeRuneInString(value)
	return unicode.IsUpper(first)
}

func isKnownScalarAttribute(value string) bool {
	value = normalizeAttribute(value)
	if value == "" || isIgnoredAttribute(value) {
		return false
	}

	return detectAttributeType(value) != "TEXT" || knownTextAttributes[value]
}

var knownTextAttributes = map[string]bool{
	"description": true, "address": true, "isbn": true, "status": true,
	"описание": true,
	"адрес":    true,
	"инн":      true,
	"код":      true,
	"марка":    true,
	"модель":   true,
	"режим":    true,
	"статус":   true,
	"фио":      true,
}
