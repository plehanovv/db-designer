package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractRelations(document nlp.Document) []model.Relation {
	var relations []model.Relation

	for _, match := range RelationMatches(document) {
		relations = append(relations, model.Relation{
			From:        normalizeEntity(match.From.Lemma),
			To:          normalizeEntity(match.To.Lemma),
			Type:        match.Type,
			Cardinality: match.Cardinality,
		})
	}

	return relations
}

func previousEntity(tokens []nlp.Token, start int, sentence int, attributeTerms map[string]bool) (nlp.Token, bool) {
	if start > 0 && isRelativePronoun(tokens[start-1].Lemma) {
		for i := start - 2; i >= 0; i-- {
			if tokens[i].Sentence != sentence {
				break
			}
			if isContextEntityNoise(tokens, i) {
				continue
			}
			if isEntity(tokens[i]) && !attributeTerms[normalizeAttribute(tokens[i].Lemma)] {
				return tokens[i], true
			}
		}
	}

	for i := start - 1; i >= 0; i-- {
		if tokens[i].Sentence != sentence {
			break
		}
		if isContextEntityNoise(tokens, i) {
			continue
		}
		if isEntity(tokens[i]) && !attributeTerms[normalizeAttribute(tokens[i].Lemma)] {
			return tokens[i], true
		}
	}

	return nlp.Token{}, false
}

func nextEntity(tokens []nlp.Token, start int, sentence int, attributeTerms map[string]bool) (nlp.Token, bool) {
	for i := start; i < len(tokens); i++ {
		if tokens[i].Sentence != sentence {
			break
		}
		if isRelationPreposition(tokens[i].Lemma) || isRelationQuantifier(tokens[i].Lemma) {
			continue
		}
		if isEntity(tokens[i]) && !attributeTerms[normalizeAttribute(tokens[i].Lemma)] {
			return tokens[i], true
		}
	}

	return nlp.Token{}, false
}

func isRelationVerb(value string) bool {
	verbs := map[string]bool{
		"belong": true, "belongs": true, "connect": true, "connects": true,
		"store": true, "stores": true, "contain": true, "contains": true,
		"include": true, "includes": true, "relate": true, "relates": true,
		"own": true, "owns": true, "place": true, "places": true,
		"create": true, "creates": true, "make": true, "makes": true,
		"buy": true, "buys": true, "borrow": true, "borrows": true,
		"pay": true, "pays": true, "use": true, "uses": true,
		"assign": true, "assigns": true,
		"work": true, "works": true,
		"enroll": true, "enrolls": true, "teach": true, "teaches": true,
		"manage": true, "manages": true, "perform": true, "performs": true,
		"execute": true, "executes": true, "participate": true, "participates": true,
		"register": true, "registers": true, "book": true, "books": true,
		"reserve": true, "reserves": true, "deliver": true, "delivers": true,
		"ship": true, "ships": true, "receive": true, "receives": true,
		"supply": true, "supplies": true, "issue": true, "issues": true,
		"sell": true, "sells": true, "visit": true, "visits": true,
		"attend": true, "attends": true, "publish": true, "publishes": true,
		"write": true, "writes": true, "sign": true, "signs": true,
		"approve": true, "approves": true, "review": true, "reviews": true,
		"open": true, "opens": true, "respond": true, "responds": true, "apply": true, "applies": true,
		"submit": true, "submits": true, "pass": true, "passes": true, "take": true, "takes": true,
		"evaluate": true, "evaluates": true, "agree": true, "agrees": true,
		"включать":      true,
		"включает":      true,
		"входит":        true,
		"входят":        true,
		"содержать":     true,
		"содержит":      true,
		"принадлежать":  true,
		"принадлежит":   true,
		"относиться":    true,
		"относится":     true,
		"связывать":     true,
		"связать":       true,
		"связан":        true,
		"связана":       true,
		"связано":       true,
		"связанный":     true,
		"использует":    true,
		"используют":    true,
		"использовать":  true,
		"назначен":      true,
		"назначена":     true,
		"хранить":       true,
		"храниться":     true,
		"хранит":        true,
		"хранится":      true,
		"хранятся":      true,
		"размещается":   true,
		"размещаться":   true,
		"оформлять":     true,
		"оформляет":     true,
		"оформляют":     true,
		"создавать":     true,
		"заказывает":    true,
		"заказывают":    true,
		"заказать":      true,
		"покупает":      true,
		"покупают":      true,
		"покупать":      true,
		"оплачивает":    true,
		"оплачивают":    true,
		"оплатить":      true,
		"создает":       true,
		"создают":       true,
		"обрабатывает":  true,
		"обрабатывают":  true,
		"арендует":      true,
		"арендуют":      true,
		"арендовать":    true,
		"работать":      true,
		"работает":      true,
		"работают":      true,
		"вести":         true,
		"ведет":         true,
		"ведут":         true,
		"записываться":  true,
		"записывается":  true,
		"управляет":     true,
		"управляют":     true,
		"выполняет":     true,
		"выполняют":     true,
		"выполнять":     true,
		"выполняется":   true,
		"выполняться":   true,
		"участвует":     true,
		"участвуют":     true,
		"регистрирует":  true,
		"бронирует":     true,
		"бронируют":     true,
		"бронировать":   true,
		"берет":         true,
		"берут":         true,
		"брать":         true,
		"резервирует":   true,
		"резервировать": true,
		"доставляет":    true,
		"доставляют":    true,
		"доставлять":    true,
		"получает":      true,
		"получают":      true,
		"получать":      true,
		"поставляет":    true,
		"поставляют":    true,
		"поставлять":    true,
		"выдает":        true,
		"выдают":        true,
		"выдавать":      true,
		"продает":       true,
		"продают":       true,
		"продавать":     true,
		"посещает":      true,
		"посещают":      true,
		"посещать":      true,
		"публикует":     true,
		"публикуют":     true,
		"публиковать":   true,
		"пишет":         true,
		"пишут":         true,
		"писать":        true,
		"подписывает":   true,
		"подписывать":   true,
		"утверждает":    true,
		"утверждать":    true,
		"проверяет":     true,
		"проверять":     true,
		"открывает":     true,
		"открывают":     true,
		"открывать":     true,
		"откликается":   true,
		"откликаются":   true,
		"откликаться":   true,
		"подает":        true,
		"подают":        true,
		"подавать":      true,
		"проходит":      true,
		"проходят":      true,
		"проходить":     true,
		"сдает":         true,
		"сдают":         true,
		"сдавать":       true,
		"оценивает":     true,
		"оценивают":     true,
		"оценивать":     true,
		"согласует":     true,
		"согласуют":     true,
		"согласовывать": true,
		"страхует":      true,
		"страхуют":      true,
		"страховать":    true,
		"проводит":      true,
		"проводят":      true,
		"ремонтирует":   true,
		"ремонтируют":   true,
		"обслуживает":   true,
		"обслуживают":   true,
		"возникает":     true,
		"возникают":     true,
		"возникать":     true,
	}

	return verbs[normalizeWord(value)]
}

func isRelationPreposition(value string) bool {
	prepositions := map[string]bool{
		"to": true, "with": true, "by": true, "in": true, "of": true,
		"к": true, "с": true, "в": true, "у": true,
		"на": true, "по": true, "во": true, "со": true,
	}

	return prepositions[normalizeWord(value)]
}

func isRelationQuantifier(value string) bool {
	quantifiers := map[string]bool{
		"many": true, "multiple": true, "several": true,
		"много": true, "несколько": true,
	}

	return quantifiers[normalizeWord(value)]
}

func relationType(value string) string {
	switch normalizeWord(value) {
	case "belong", "belongs", "work", "works", "принадлежать", "принадлежит", "относиться", "относится", "работает", "работают", "храниться", "хранится", "размещается", "размещаться", "выполняется", "выполняться":
		return "belongs_to"
	case "contain", "contains", "include", "includes", "содержать", "содержит", "включать", "включает", "входит", "входят":
		return "contains"
	case "connect", "connects", "relate", "relates", "use", "uses", "assign", "assigns", "enroll", "enrolls", "teach", "teaches",
		"manage", "manages", "perform", "performs", "execute", "executes", "participate", "participates", "respond", "responds", "apply", "applies",
		"visit", "visits", "attend", "attends", "review", "reviews", "approve", "approves", "pass", "passes", "take", "takes", "evaluate", "evaluates", "agree", "agrees",
		"связывать", "связать", "связан", "связана", "связано", "связанный", "использует", "используют", "использовать", "назначен", "назначена",
		"управляет", "управляют", "управлять", "выполняет", "выполняют", "выполнять", "участвует", "участвуют", "участвовать", "посещает", "посещают", "посещать", "утверждает", "утверждать", "проверяет", "проверять", "откликается", "откликаются", "откликаться", "проходит", "проходят", "проходить", "сдает", "сдают", "сдавать", "оценивает", "оценивают", "оценивать", "согласует", "согласуют", "согласовывать", "проводит", "проводят", "проводить", "ремонтирует", "ремонтируют", "ремонтировать", "обслуживает", "обслуживают", "обслуживать", "возникает", "возникают", "возникать":
		return "associated_with"
	case "place", "places", "create", "creates", "make", "makes", "buy", "buys", "borrow", "borrows", "pay", "pays", "open", "opens", "submit", "submits", "register", "registers", "book", "books", "reserve", "reserves", "deliver", "delivers", "ship", "ships", "receive", "receives", "supply", "supplies", "issue", "issues", "sell", "sells", "publish", "publishes", "write", "writes", "sign", "signs",
		"оформлять", "оформляет", "оформляют", "заказывает", "заказывают", "заказать", "покупает", "покупают", "покупать", "оплачивает", "оплачивают", "оплатить", "создает", "создают", "создавать", "обрабатывает", "обрабатывают", "обрабатывать", "арендует", "арендуют", "арендовать", "вести", "ведет", "ведут",
		"открывает", "открывают", "открывать", "подает", "подают", "подавать", "страхует", "страхуют", "страховать", "регистрирует", "регистрировать", "бронирует", "бронируют", "бронировать", "берет", "берут", "брать", "резервирует", "резервировать", "доставляет", "доставляют", "доставлять", "получает", "получают", "получать", "поставляет", "поставляют", "поставлять", "выдает", "выдают", "выдавать", "продает", "продают", "продавать", "публикует", "публикуют", "публиковать", "пишет", "пишут", "писать", "подписывает", "подписывать":
		return "has"
	default:
		return "associated_with"
	}
}

func relationCardinality(value string) string {
	if isManyToManyRelationVerb(value) {
		return "many-to-many"
	}
	if isLinkedManyToOneRelationVerb(value) {
		return "many-to-one"
	}
	if isActionOneToManyRelationVerb(value) {
		return "one-to-many"
	}

	switch relationType(value) {
	case "belongs_to":
		return "many-to-one"
	case "contains", "has":
		return "one-to-many"
	default:
		return "unspecified"
	}
}

func isManyToManyRelationVerb(value string) bool {
	verbs := map[string]bool{
		"enroll": true, "enrolls": true, "teach": true, "teaches": true,
		"pass": true, "passes": true, "take": true, "takes": true, "borrow": true, "borrows": true,
		"вести":        true,
		"ведет":        true,
		"ведут":        true,
		"записываться": true,
		"записывается": true,
		"сдавать":      true,
		"сдает":        true,
		"сдают":        true,
		"берет":        true,
		"берут":        true,
		"брать":        true,
	}

	return verbs[normalizeWord(value)]
}

func isLinkedManyToOneRelationVerb(value string) bool {
	verbs := map[string]bool{
		"belong": true, "belongs": true, "connect": true, "connects": true,
		"relate": true, "relates": true, "assign": true, "assigns": true,
		"принадлежать": true,
		"принадлежит":  true,
		"относиться":   true,
		"относится":    true,
		"связывать":    true,
		"связать":      true,
		"связан":       true,
		"связана":      true,
		"связано":      true,
		"связанный":    true,
		"назначен":     true,
		"назначена":    true,
		"использует":   true,
		"используют":   true,
		"использовать": true,
		"храниться":    true,
		"хранится":     true,
		"размещается":  true,
		"размещаться":  true,
		"выполняется":  true,
		"выполняться":  true,
	}

	return verbs[normalizeWord(value)]
}

func isActionOneToManyRelationVerb(value string) bool {
	verbs := map[string]bool{
		"evaluate": true, "evaluates": true, "review": true, "reviews": true,
		"approve": true, "approves": true,
		"оценивать":  true,
		"оценивает":  true,
		"оценивают":  true,
		"проверяет":  true,
		"проверять":  true,
		"утверждает": true,
		"утверждать": true,
		"выполняет":  true,
		"выполняют":  true,
		"выполнять":  true,
		"возникает":  true,
		"возникают":  true,
		"возникать":  true,
	}

	return verbs[normalizeWord(value)]
}
