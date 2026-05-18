package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractRelations(document nlp.Document) []model.Relation {

	var relations []model.Relation

	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {

		first := tokens[i]
		second := tokens[i+1]
		third := tokens[i+2]

		if isEntity(first) &&
			isRelationVerb(second.Lemma) &&
			isEntity(third) {

			relations = append(relations, model.Relation{
				From: normalizeEntity(first.Lemma),
				To:   normalizeEntity(third.Lemma),
				Type: second.Lemma,
			})
		}
	}

	return relations
}

func isRelationVerb(value string) bool {

	verbs := map[string]bool{
		"have":    true,
		"contain": true,
		"include": true,
		"belong":  true,
		"connect": true,
		"store":   true,
	}

	return verbs[value]
}
