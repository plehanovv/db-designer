package analyzer

import "strings"

func detectAttributeType(attribute string) string {
	attribute = normalizeAttribute(attribute)
	if strings.HasPrefix(attribute, "дата_") {
		return "DATE"
	}

	switch attribute {

	case "age", "year", "floor", "quantity", "count", "stock", "duration", "rating", "score", "capacity", "seats", "places", "attempt", "priority",
		"возраст", "год", "этаж", "количество", "номер",
		"остаток", "длительность", "рейтинг", "балл", "оценка", "место", "приоритет", "гарантия":
		return "INTEGER"

	case "price", "cost", "amount", "discount", "salary", "budget", "balance", "tax", "total", "weight", "volume", "rate", "fee", "payment",
		"скидка", "цена", "стоимость", "сумма",
		"зарплата", "оклад", "бюджет", "баланс", "налог", "итог", "вес", "объем", "площадь", "тариф", "комиссия", "платеж", "премия", "выплата":
		return "NUMERIC(12,2)"

	case "date", "created", "updated", "deadline", "birthdate", "дата", "дедлайн", "срок":
		return "DATE"

	case "time", "start_time", "end_time", "время":
		return "TIME"

	case "email", "почта":
		return "VARCHAR(255)"

	case "phone", "телефон":
		return "VARCHAR(20)"

	case "name", "title", "login", "code", "author", "group", "position", "role", "status", "model", "brand", "passport", "vin",
		"address", "city", "country", "region", "street", "category", "type", "kind", "level", "description", "comment", "note", "url", "slug",
		"автор", "группа", "должность", "имя", "код", "марка", "модель", "название", "логин", "паспорт", "роль", "статус", "тема", "госномер", "специальность",
		"адрес", "город", "страна", "регион", "улица", "категория", "тип", "уровень", "комментарий", "примечание", "ссылка", "класс", "кабинет", "маршрут", "дозировка", "жанр", "серия":
		return "VARCHAR(255)"

	default:
		return "TEXT"
	}
}
