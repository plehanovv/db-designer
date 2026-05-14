package analyzer

func detectAttributeType(attribute string) string {

	switch attribute {

	case "age":
		return "INTEGER"

	case "email":
		return "VARCHAR(255)"

	case "phone":
		return "VARCHAR(20)"

	default:
		return "TEXT"
	}
}
