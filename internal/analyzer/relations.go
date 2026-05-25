package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractRelations(document nlp.Document) []model.Relation {

	var relations []model.Relation
	seen := make(map[string]bool)

	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {

		first := tokens[i]
		second := tokens[i+1]
		third := tokens[i+2]

		if isEntity(first) &&
			isRelationVerb(second.Lemma) &&
			isEntity(third) {
			from := normalizeEntity(first.Lemma)
			to := normalizeEntity(third.Lemma)
			key := from + "->" + to
			if seen[key] {
				continue
			}

			relations = append(relations, model.Relation{
				From:        from,
				To:          to,
				Type:        relationType(second.Lemma),
				Cardinality: relationCardinality(second.Lemma),
			})
			seen[key] = true
		}
	}

	return relations
}

func isRelationVerb(value string) bool {

	verbs := map[string]bool{
		"have": true, "has": true, "contain": true, "contains": true,
		"include": true, "includes": true, "belong": true, "belongs": true,
		"connect": true, "connects": true, "store": true, "stores": true,
		"\u0438\u043c\u0435\u0442\u044c":                                           true,
		"\u0438\u043c\u0435\u0435\u0442":                                           true,
		"\u0438\u043c\u0435\u044e\u0442":                                           true,
		"\u0441\u043e\u0434\u0435\u0440\u0436\u0430\u0442\u044c":                   true,
		"\u0441\u043e\u0434\u0435\u0440\u0436\u0438\u0442":                         true,
		"\u0432\u043a\u043b\u044e\u0447\u0430\u0442\u044c":                         true,
		"\u0432\u043a\u043b\u044e\u0447\u0430\u0435\u0442":                         true,
		"\u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0430\u0442\u044c": true,
		"\u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0438\u0442":       true,
		"\u0441\u0432\u044f\u0437\u044b\u0432\u0430\u0442\u044c":                   true,
		"\u0441\u0432\u044f\u0437\u0430\u043d":                                     true,
		"\u0441\u0432\u044f\u0437\u0430\u043d\u0430":                               true,
		"\u0445\u0440\u0430\u043d\u0438\u0442\u044c":                               true,
		"\u0445\u0440\u0430\u043d\u0438\u0442":                                     true,
	}

	return verbs[normalizeWord(value)]
}

func relationType(value string) string {
	switch normalizeWord(value) {
	case "belong", "belongs", "\u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0430\u0442\u044c", "\u043f\u0440\u0438\u043d\u0430\u0434\u043b\u0435\u0436\u0438\u0442":
		return "belongs_to"
	case "contain", "contains", "include", "includes", "\u0441\u043e\u0434\u0435\u0440\u0436\u0430\u0442\u044c", "\u0441\u043e\u0434\u0435\u0440\u0436\u0438\u0442", "\u0432\u043a\u043b\u044e\u0447\u0430\u0442\u044c", "\u0432\u043a\u043b\u044e\u0447\u0430\u0435\u0442":
		return "contains"
	case "connect", "connects", "\u0441\u0432\u044f\u0437\u044b\u0432\u0430\u0442\u044c", "\u0441\u0432\u044f\u0437\u0430\u043d", "\u0441\u0432\u044f\u0437\u0430\u043d\u0430":
		return "associated_with"
	default:
		return "has"
	}
}

func relationCardinality(value string) string {
	switch relationType(value) {
	case "belongs_to":
		return "many-to-one"
	case "contains", "has":
		return "one-to-many"
	default:
		return "unspecified"
	}
}
