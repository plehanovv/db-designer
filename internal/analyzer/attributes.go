package analyzer

import (
	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/nlp"
)

func ExtractAttributes(
	document nlp.Document,
	entityMap map[string]*model.Entity,
) {

	tokens := document.Tokens

	for i := 0; i < len(tokens)-2; i++ {

		first := tokens[i]
		second := tokens[i+1]
		third := tokens[i+2]

		if isEntity(first) &&
			second.Lemma == "have" {

			entityName := normalizeEntity(first.Lemma)

			entity := entityMap[entityName]

			if entity == nil {
				continue
			}

			entity.Attributes = append(
				entity.Attributes,
				model.Attribute{
					Name: third.Lemma,
					Type: detectAttributeType(third.Lemma),
				},
			)
		}
	}
}
