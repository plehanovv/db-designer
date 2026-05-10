package analyzer

import "db-designer-vkr/internal/model"

func ExtractRelations(words []string) []model.Relation {

	var relations []model.Relation

	for i := 0; i < len(words)-2; i++ {

		first := clean(words[i])
		second := clean(words[i+1])
		third := clean(words[i+2])

		if isEntity(first) && isVerb(second) && isEntity(third) {

			relations = append(relations, model.Relation{
				From: first,
				To:   third,
				Type: second,
			})
		}
	}

	return relations
}

func isVerb(word string) bool {

	verbs := []string{
		"enrolls",
		"teaches",
		"belongs",
		"uses",
		"manages",
		"creates",
	}

	for _, v := range verbs {

		if word == v {
			return true
		}
	}

	return false
}
