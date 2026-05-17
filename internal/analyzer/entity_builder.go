package analyzer

import "db-designer-vkr/internal/model"

func BuildEntities(
	candidates []model.EntityCandidate,
) map[string]*model.Entity {

	entities := make(map[string]*model.Entity)

	for _, candidate := range candidates {

		entities[candidate.Name] = &model.Entity{
			Name:       candidate.Name,
			Attributes: []model.Attribute{},
		}
	}

	return entities
}
