package analyzer

import (
	"db-designer-vkr/internal/knowledge"
	"db-designer-vkr/internal/model"
)

func ExtractAttributes(
	words []string,
	entityMap map[string]*model.Entity,
) {

	for i := 0; i < len(words)-2; i++ {

		first := words[i]
		second := words[i+1]
		third := words[i+2]

		if knowledge.Dictionary[first] == knowledge.EntityType &&
			second == "has" &&
			knowledge.Dictionary[third] == knowledge.AttributeType {

			entity := entityMap[first]

			entity.Attributes = append(
				entity.Attributes,
				model.Attribute{
					Name: third,
					Type: detectAttributeType(third),
				},
			)
		}
	}
}
