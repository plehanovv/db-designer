package analyzer

func detectAttributeType(attribute string) string {
	attribute = normalizeAttribute(attribute)

	switch attribute {

	case "age", "\u0432\u043e\u0437\u0440\u0430\u0441\u0442", "\u043a\u043e\u043b\u0438\u0447\u0435\u0441\u0442\u0432\u043e", "\u043d\u043e\u043c\u0435\u0440":
		return "INTEGER"

	case "price", "cost", "amount", "\u0446\u0435\u043d\u0430", "\u0441\u0442\u043e\u0438\u043c\u043e\u0441\u0442\u044c", "\u0441\u0443\u043c\u043c\u0430":
		return "NUMERIC(12,2)"

	case "date", "created", "updated", "\u0434\u0430\u0442\u0430":
		return "DATE"

	case "email", "\u043f\u043e\u0447\u0442\u0430":
		return "VARCHAR(255)"

	case "phone", "\u0442\u0435\u043b\u0435\u0444\u043e\u043d":
		return "VARCHAR(20)"

	case "name", "title", "login", "\u0438\u043c\u044f", "\u043d\u0430\u0437\u0432\u0430\u043d\u0438\u0435", "\u043b\u043e\u0433\u0438\u043d":
		return "VARCHAR(255)"

	default:
		return "TEXT"
	}
}
