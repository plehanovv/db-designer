package analyzer

import (
	"db-designer-vkr/internal/knowledge"
	"db-designer-vkr/internal/model"
)

func ExtractRelations(words []string) []model.Relation {
	var relations []model.Relation

	for i := 0; i < len(words)-2; i++ {
		first := words[i]
		second := words[i+1]
		third := words[i+2]

		if knowledge.Dictionary[first] == knowledge.EntityType &&
			knowledge.Dictionary[second] == knowledge.VerbType &&
			knowledge.Dictionary[third] == knowledge.EntityType {

			relations = append(relations, model.Relation{
				From: first,
				To:   third,
				Type: second,
			})
		}
	}

	return relations
}
