package analyzer

import "db-designer-vkr/internal/model"

func ExtractAttributes(
	words []string,
	entityMap map[string]*model.Entity,
) {

	for i := 0; i < len(words)-2; i++ {

		first := clean(words[i])
		second := clean(words[i+1])
		third := clean(words[i+2])

		if isEntity(first) && second == "has" && !isEntity(third) {

			entity := entityMap[first]

			entity.Attributes = append(
				entity.Attributes,
				model.Attribute{
					Name: third,
					Type: "TEXT",
				},
			)
		}
	}
}
